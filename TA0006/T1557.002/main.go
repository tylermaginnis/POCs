package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"net"
)

func main() {
	for {
		// Execute the arp -a command to list ARP entries
		arpCmd := exec.Command("arp", "-a")
		output, err := arpCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error executing ARP command:", err)
			return
		}

		// Parse ARP table and extract IP-MAC pairs
		arpTable := parseARPTable(string(output))

		// Select IP address and corresponding MAC from ARP table
		selectedIP, selectedMAC := selectIPAddressAndMAC(arpTable)

		localIP, err := getLocalIPWithMatchingSubnet(selectedIP)


		if selectedIP != "" {
			// Announce the selected IP address with the corresponding MAC
			fmt.Printf("Announcing IP address: %s, MAC address: %s\n", selectedIP, selectedMAC)
			announceIP(localIP, selectedMAC)
		}

		fmt.Println("Press Enter to announce another IP address, or type 'exit' to quit...")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
	}
}

// Parse ARP table and extract IP-MAC pairs
func parseARPTable(output string) map[string]string {
	arpTable := make(map[string]string)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			ipAddr := fields[0]
			macAddr := fields[1]
			arpTable[ipAddr] = macAddr
		}
	}
	return arpTable
}

// Select IP address and corresponding MAC from ARP table
func selectIPAddressAndMAC(arpTable map[string]string) (string, string) {
	fmt.Println("Available devices:")
	i := 1
	var selectedIP, selectedMAC string
	for ip, mac := range arpTable {
		fmt.Printf("%d - IP: %s, MAC: %s\n", i, ip, mac)
		i++
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number of the device to announce (or press Enter to skip): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		deviceNum, err := strconv.Atoi(input)
		if err != nil || deviceNum < 1 || deviceNum > i {
			fmt.Println("Invalid device number.")
		} else {
			i = 1
			for ip, mac := range arpTable {
				if i == deviceNum {
					selectedIP = ip
					selectedMAC = mac
					break
				}
				i++
			}
		}
	}

	return selectedIP, selectedMAC
}

// Send ARP announcement for the selected IP address in an endless loop
func announceIP(ip, mac string) {
	for {
		cmd := exec.Command("arp", "-s", ip, mac)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error announcing IP address %s: %v\n", ip, err)
		} else {
			fmt.Printf("Announced IP address: %s, MAC address: %s\n", ip, mac)
		}

		// Sleep for a few seconds (you can adjust the duration)
		time.Sleep(5 * time.Second)
	}
}

func getLocalIPWithMatchingSubnet(selectedIP string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	selectedIPParts := strings.Split(selectedIP, ".")
	if len(selectedIPParts) != 4 {
		return "", fmt.Errorf("Invalid selected IP address format")
	}

	selectedSubnetMask := selectedIPParts[0] + "." + selectedIPParts[1] + "." + selectedIPParts[2] + ".0"

	for _, intf := range interfaces {
		addrs, err := intf.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ipParts := strings.Split(ipNet.IP.String(), ".")
				if len(ipParts) != 4 {
					continue
				}

				localSubnetMask := ipParts[0] + "." + ipParts[1] + "." + ipParts[2] + ".0"
				if localSubnetMask == selectedSubnetMask {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("Matching local IP not found")
}
