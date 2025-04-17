package main

import "net"

type InterfaceConfig struct {
	Name            string
	Description     string
	Physical        bool
	IPAddress       net.IP
	NetworkMask     net.IP
	NetworkMaskSize int
	NetworkAddress  net.IP
}

type SwitchConfig struct {
	Hostname       string
	SwitchID       int
	Interfaces     map[string]InterfaceConfig
	InterfaceNames []string
}
