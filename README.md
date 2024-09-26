# Poweredge Fan Control
Manually sets fan speeds so your poweredge doesn't sound like it's going to take off.

*Use this tool at your own risk*

## Requirements
- iDRAC version must be less than `3.30.30.30`
- IPMI must be enabled
	- In the iDRAC sidebar: `Overview > iDRAC Settings > Network > IPMI Settings`
- `ipmitool` installed
- `nvidia-smi` installed

## Layout
- `fanctl` Example implementation for my specific case
- `pkg`
	- `ipmitool` Invokes and parses the output of [ipmitool](https://git.launchpad.net/ubuntu/+source/ipmitool/)
	- `nvidia` Wraps [nvidia-smi](https://developer.nvidia.com/system-management-interface)

Credit to [spxlabs](https://www.spxlabs.com/blog/2019/3/16/silence-your-dell-poweredge-server) for helping with the setup and commands. 
