package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/3th1nk/cidr"
)

func main() {
	cfgs := ParseCiscoConfigs()
	keaCfgs := make([]KeaDHCP4SubnetConfig, 0, len(cfgs)*64)
	for _, cfg := range cfgs {
		fmt.Println("------------------------")
		fmt.Println(cfg.Hostname)
		fmt.Println("------------------------")
		for _, key := range cfg.InterfaceNames {
			fmt.Println(key)
			fmt.Print("\t")
			fmt.Println(cfg.Interfaces[key])
			fmt.Print("\t")
			fmt.Println(removeAlphabet(key))
		}
		keaCfgs = append(keaCfgs, createKeaConfig(cfg)...)
	}
	bytes, err := json.MarshalIndent(keaCfgs, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("output/kea-config.json", bytes, 0744)
	if err != nil {
		panic(err)
	}
}

func createKeaConfig(sw SwitchConfig) []KeaDHCP4SubnetConfig {
	intCfgs := make([]KeaDHCP4SubnetConfig, 0, len(sw.InterfaceNames))
	for _, inter := range sw.InterfaceNames {
		if !(sw.Interfaces[inter].Physical) {
			continue
		}
		if sw.Interfaces[inter].IPAddress.IsUnspecified() {
			continue
		}
		if sw.Interfaces[inter].NetworkMaskSize < 30 {
			continue
		}
		intCfgs = append(intCfgs, createKeaSubnetConfig(sw.SwitchID, sw.Interfaces[inter]))
	}
	return intCfgs
}

func createKeaSubnetConfig(swid int, inter InterfaceConfig) KeaDHCP4SubnetConfig {
	// Compute an ID for this interface that will not change with varying names
	intId := 0
	if compareStringStart("mgmt", inter.Name) {
		intId, _ = strconv.Atoi(string(inter.Name[len(inter.Name)-1]))
		intId += 900
		intId *= 10
	} else {
		// Handle breakouts
		segments := strings.Split(removeAlphabet(inter.Name), "/")
		if len(segments)-1 > 1 {
			val, _ := strconv.Atoi(segments[len(segments)-2])
			intId = val * 1000
		}
		if len(segments)-1 > 0 {
			val, _ := strconv.Atoi(segments[len(segments)-1])
			intId += val * 10
		}
	}
	intId = swid*10000 + intId
	pool := getHostIP(inter).String() + " - " + getHostIP(inter).String()
	fmt.Println(inter.NetworkAddress)
	return KeaDHCP4SubnetConfig{
		Subnet: inter.NetworkAddress.String() + "/" + strconv.Itoa(inter.NetworkMaskSize),
		ID:     uint32(intId),
		Pools: []KeaDHCP4SubnetPoolConfig{
			{
				Pool: pool,
			},
		},
		OptionData: []KeaDHCP4SubnetOptionConfig{
			{
				Name: ROUTERS,
				Data: inter.IPAddress.String(),
			},
		},
	}
}

func ParseCiscoConfigs() map[string]SwitchConfig {
	files, err := os.ReadDir("configs")
	if err != nil {
		panic(err)
	}
	switchConfigs := make(map[string]SwitchConfig)
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}
		data, err := os.ReadFile("configs/" + file.Name())
		if err != nil {
			panic(err)
		}
		cfg := ParseCiscoConfig(string(data))
		switchConfigs[cfg.Hostname] = cfg
	}
	return switchConfigs
}

func ParseCiscoConfig(textConfig string) SwitchConfig {
	config := SwitchConfig{
		Interfaces: make(map[string]InterfaceConfig),
	}

	cfg := strings.Split(strings.ReplaceAll(textConfig, "\r\n", "\n"), "\n")

	inIntConfig := false
	currentInt := InterfaceConfig{}
	for _, line := range cfg {
		if compareStringStart("hostname ", line) {
			config.Hostname = line[9:]
		}
		if inIntConfig && !compareStringStart("  ", line) {
			inIntConfig = false
			config.Interfaces[currentInt.Name] = currentInt
			config.InterfaceNames = append(config.InterfaceNames, currentInt.Name)
		}
		if inIntConfig {
			if compareStringStart("  description ", line) {
				currentInt.Description = line[14:]
				continue
			}
			if compareStringStart("  ip address ", line) {
				// Ignore secondary addresses
				if compareStringEnd("secondary", line) {
					continue
				}
				addr, err := cidr.Parse(line[13:])
				if err != nil {
					panic(err)
				}
				currentInt.IPAddress = addr.IP()
				currentInt.NetworkMask = addr.Mask()
				leading, _ := addr.MaskSize()
				currentInt.NetworkMaskSize = leading
				currentInt.NetworkAddress = addr.Network()
				continue
			}
		}
		if compareStringStart("interface ", line) {
			inIntConfig = true
			currentInt = InterfaceConfig{
				Name: line[10:],
			}
			if compareStringStart("Eth", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("Gi", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("Fa", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("Te", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("Tw", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("Fo", currentInt.Name) {
				currentInt.Physical = true
			}
			if compareStringStart("mgmt", currentInt.Name) {
				currentInt.Physical = true
			}
			continue
		}
	}

	// Get switch ID from loopback0 fourth octet
	config.SwitchID = int(config.Interfaces["loopback0"].IPAddress.To4()[3])

	return config
}

func compareStringStart(search string, str string) bool {
	if len(str) < len(search) {
		return false
	}
	return strings.Compare(search, str[0:len(search)]) == 0
}

func compareStringEnd(search string, str string) bool {
	if len(str) < len(search) {
		return false
	}
	return strings.Compare(search, str[len(str)-len(search):]) == 0
}

// Returns the free (host) ip address for a /30 or /31 link
func getHostIP(inter InterfaceConfig) net.IP {
	// Deal with a /31
	if inter.NetworkMaskSize == 31 {
		if inter.IPAddress.Equal(inter.NetworkAddress) {
			deviceIP := make(net.IP, len(inter.NetworkAddress))
			copy(deviceIP, inter.NetworkAddress)
			deviceIP = deviceIP.To4()
			deviceIP[3] += 1
			return deviceIP
		} else {
			return inter.NetworkAddress
		}
	}
	// Deal with a /30
	if inter.NetworkMaskSize == 30 {
		ip1 := make(net.IP, len(inter.NetworkAddress))
		copy(ip1, inter.NetworkAddress)
		ip1.To4()
		ip1[3] += 1
		ip2 := make(net.IP, len(inter.NetworkAddress))
		copy(ip2, inter.NetworkAddress)
		ip2.To4()
		ip2[3] += 2
		if inter.IPAddress.Equal(ip1) {
			return ip2
		} else {
			return ip1
		}
	}
	return nil
}

// Leaves only numbers and other non a-z,A-Z characters
func removeAlphabet(str string) string {
	return regexp.MustCompile("[A-z-]").ReplaceAllString(str, "")
}
