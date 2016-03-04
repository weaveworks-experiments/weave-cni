package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/appc/cni/pkg/skel"
	"github.com/appc/cni/pkg/types"
)

func main() {
	skel.PluginMain(cmdAdd, cmdDel)
}

func cmdAdd(args *skel.CmdArgs) error {
	// extract the things we care about
	conf, err := loadIPAMConf(args.StdinData)
	if err != nil {
		return err
	}

	containerID := args.ContainerID
	var url string

	if conf.Subnet == "" {
		url = fmt.Sprintf("/ip/%s", containerID)
	} else {
		_, err = types.ParseCIDR(conf.Subnet)
		if err != nil {
			return fmt.Errorf("subnet given in config, but not parseable: %s", err)
		}
		url = fmt.Sprintf("/ip/%s/%s", containerID, conf.Subnet)
	}

	var ip string
	if ip, err = httpPost(url); err != nil {
		return err
	}
	ipnet, err := types.ParseCIDR(ip)
	if err != nil {
		return err
	}
	result := types.Result{
		IP4: &types.IPConfig{
			IP:      *ipnet,
			Gateway: conf.Gateway,
		},
	}
	return result.Print()
}

func cmdDel(args *skel.CmdArgs) error {
	containerID := args.ContainerID
	url := fmt.Sprintf("/ip/%s", containerID)
	if err := httpDelete(url); err != nil {
		return err
	}
	result := types.Result{}
	return result.Print()
}

type ipamConf struct {
	Subnet  string `json:"subnet,omitempty"`
	Gateway net.IP `json:"gateway,omitempty"`
}

type netConf struct {
	IPAM *ipamConf `json:"ipam"`
}

func loadIPAMConf(stdinData []byte) (*ipamConf, error) {
	var conf netConf
	return conf.IPAM, json.Unmarshal(stdinData, &conf)
}

func httpPost(path string) (string, error) {
	resp, err := http.Post("http://localhost:6784"+path, "", nil)
	if err != nil {
		return "", err
	}
	if !okResponse(resp) {
		return "", fmt.Errorf("non-OK status code %d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func httpDelete(path string) error {
	req, err := http.NewRequest("DELETE", "http://localhost:6784"+path, nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if !okResponse(resp) {
		return fmt.Errorf("non-OK status code %d", resp.StatusCode)
	}
	return nil
}

func okResponse(response *http.Response) bool {
	switch response.StatusCode {
	case http.StatusOK:
		return true
	case http.StatusCreated:
		return true
	case http.StatusAccepted:
		return true
	default:
		return false
	}
}
