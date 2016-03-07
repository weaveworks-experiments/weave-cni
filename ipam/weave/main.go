package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/appc/cni/pkg/skel"
	"github.com/appc/cni/pkg/types"
	"github.com/weaveworks/weave/api"
)

func main() {
	skel.PluginMain(cmdAdd, cmdDel)
}

var (
	log      = logrus.New()
	weaveapi = api.NewClient("localhost:6784", log)
)

func cmdAdd(args *skel.CmdArgs) error {
	// extract the things we care about
	conf, err := loadIPAMConf(args.StdinData)
	if err != nil {
		return err
	}

	containerID := args.ContainerID
	var ipnet *net.IPNet

	if conf.Subnet == "" {
		ipnet, err = weaveapi.AllocateIP(containerID)
	} else {
		var subnet *net.IPNet
		subnet, err = types.ParseCIDR(conf.Subnet)
		if err != nil {
			return fmt.Errorf("subnet given in config, but not parseable: %s", err)
		}
		ipnet, err = weaveapi.AllocateIPInSubnet(containerID, subnet)
	}

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
	if err := weaveapi.ReleaseIP(containerID); err != nil {
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
