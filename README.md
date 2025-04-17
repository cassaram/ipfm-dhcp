# IPFM-DHCP

IPFM-DHCP is a tool designed to extract IP address information from interfaces in network switch config files and derrive configurations for a DHCP server. This was primarily designed for the context of an IP Fabric for Media (Cisco term for 2110-style spine and leaf, multicast-heavy topologies) in order to minimize the time needed reconfiguring device IP addressing when moving them between network ports. Supports both /30 and /31 addressing to host devices. Outputs DHCP server configuration steps to a file, allowing users to double check and verify the output before comitting it, as the output may contain network links that should not be assigned by DHCP.

## Currently Supported Switches

Cisco IOS-XE
Cisco IOS-XR
Cisco NX-OS

## Currently Supported DHCP Servers

ISC-Kea


## Implementation Notes

ISC-Kea requires an ID to be generated per subnet. These IDs are what tie leases to subnets, so they should be statically defined to ensure that a config reload does not tie leases to the wrong subnet.

This program assumes that each switch has a unique ID as the 4th Octet of the loopback0 interface, which is used to derrive the IDs.

## How to use

1. Place switch configuration files into the ./configs folder
2. Run the program
3. Result partial config files for DHCP servers are written to the ./output folder