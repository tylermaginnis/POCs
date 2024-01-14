package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type NetworkInterfaceInfo struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
	Hostnames []string `json:"hostnames,omitempty"`
}

type ArpEntry struct {
	IPAddress  string   `json:"ip_address"`
	MACAddress string   `json:"mac_address"`
	Hostnames  []string `json:"hostnames,omitempty"`
}

func main() {
	networkInfo, err := gatherNetworkInfo()
	if err != nil {
		fmt.Println("Error gathering network info:", err)
		os.Exit(1)
	}

	arpEntries, err := getArpCache()
	if err != nil {
		fmt.Println("Error gathering ARP cache:", err)
		os.Exit(1)
	}

	output := struct {
		NetworkInterfaces []NetworkInterfaceInfo `json:"network_interfaces"`
		ArpCache          []ArpEntry             `json:"arp_cache"`
	}{
		NetworkInterfaces: networkInfo,
		ArpCache:          arpEntries,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput))
}

func gatherNetworkInfo() ([]NetworkInterfaceInfo, error) {
	var networkInfo []NetworkInterfaceInfo

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		var ifaceInfo NetworkInterfaceInfo
		ifaceInfo.Name = iface.Name

		addrs, err := iface.Addrs()
		if err != nil {
			continue // Skipping interface if there's an error getting addresses
		}

		for _, addr := range addrs {
			ifaceInfo.Addresses = append(ifaceInfo.Addresses, addr.String())

			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue // Skipping address if there's an error parsing it
			}

			hostnames, err := net.LookupAddr(ip.String())
			if err == nil {
				ifaceInfo.Hostnames = append(ifaceInfo.Hostnames, hostnames...)
			}
		}

		networkInfo = append(networkInfo, ifaceInfo)
	}

	return networkInfo, nil
}

func getArpCache() ([]ArpEntry, error) {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseArpOutput(string(output))
}

func parseArpOutput(output string) ([]ArpEntry, error) {
	var entries []ArpEntry
	scanner := bufio.NewScanner(strings.NewReader(output))
	arpRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s+([0-9a-f-]+)\s`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := arpRegex.FindStringSubmatch(line)
		if matches != nil && len(matches) > 2 {
			ip := matches[1]
			mac := matches[2]

			hostnames, _ := net.LookupAddr(ip)

			entry := ArpEntry{
				IPAddress:  ip,
				MACAddress: mac,
				Hostnames:  hostnames,
			}
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
