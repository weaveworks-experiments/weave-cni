## Using Weave with CNI

Assuming you have the plugin binary `weave-ipam` (see below if you
want to build it), put that binary somewhere on your
[`CNI_PATH`][cni]. Then, make a CNI network configuration using the
bridge plugin as the network type, and the `weave-ipam` plugin as the
IPAM type. Here's an example:

```
{
    "name": "weave",
    "type": "bridge",
    "bridge": "weave",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "weave-ipam",
        "subnet": "10.32.127.0/24",
        "gateway": "10.32.0.1",
        "routes": [
            { "dst": "0.0.0.0/0" }
        ]
    }
}
```

You should then be able to use `cnitool`, or other libCNI
applications, to give containers as weave interface.

## Building the plugin

The build assumes you are running it on Linux. In the cloned repository:

    weave-cni $ make

This deposits the plugin binary `weave-ipam` in the top directory.

[cni]: https://github.com/appc/cni#included-plugins
