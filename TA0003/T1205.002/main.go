package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	snapshotLen  int32         = 1024
	promiscuous  bool          = false
	err          error
	timeout      time.Duration = pcap.BlockForever
	handle       *pcap.Handle
	targetPort   layers.TCPPort = 80 // Listening on port 80
)

func main() {
	fmt.Println("Fetching available network interfaces...")

	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error finding devices:", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No devices found. Make sure you have the necessary permissions.")
		return
	}

	// Display available devices
	for i, device := range devices {
		fmt.Printf("%d. %s - %s\n", i+1, device.Name, device.Description)
	}

	// User input to select device
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the number of the device to use for packet capturing: ")
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > len(devices) {
		fmt.Println("Invalid input. Exiting.")
		return
	}

	selectedDevice := devices[choice-1]
	fmt.Printf("You selected: %s\n", selectedDevice.Name)
	fmt.Printf("Listening for custom packets on port %d\n", targetPort)

	// Open the device for capturing
	handle, err = pcap.OpenLive(selectedDevice.Name, snapshotLen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device %s: %v\n", selectedDevice.Name, err)
		return
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			if tcp.DstPort == targetPort {
				fmt.Printf("Received TCP packet from %s to %s on port %d\n", packet.NetworkLayer().NetworkFlow().Src(), packet.NetworkLayer().NetworkFlow().Dst(), tcp.DstPort)
				if isCustomPacket(packet) {
					fmt.Println("Custom packet detected. Taking action...")
					takeAction()
				}
			}
		}
	}
}

func isCustomPacket(packet gopacket.Packet) bool {
	// Check for the application layer in the packet
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		// Convert the payload to a string for easier handling
		payload := string(applicationLayer.Payload())

		// Check if the payload contains the specific string
		if strings.Contains(payload, "custom_packet") {
			// Additionally, check for standard HTTP headers
			if strings.Contains(payload, "HTTP/1.1") || strings.Contains(payload, "HTTP/2.0") {
				fmt.Println("Custom HTTP packet with 'custom_packet' string detected.")
				return true
			}
		}
	}

	// If the packet does not contain the specific string or standard HTTP headers
	return false
}


func takeAction() {
	// Implement the action to be taken after detecting the custom packet
	fmt.Println("Action taken: Custom packet detected.")
	// Define actions here, such as triggering a service, sending a response, etc.
}