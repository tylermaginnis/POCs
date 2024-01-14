package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	snapshotLen        int32 = 1024
	promiscuous         bool  = false
	err                 error
	timeout                   = pcap.BlockForever
	handle              *pcap.Handle
	knockSequence            = []int{1234, 5678, 9012} // Example ports
	currentKnockIndex        = 0
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
	fmt.Printf("Waiting for the knock sequence... 1234, 5678, 9012\n")

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

			fmt.Printf("Received TCP packet from %s to %s on port %d\n", packet.NetworkLayer().NetworkFlow().Src(), packet.NetworkLayer().NetworkFlow().Dst(), tcp.DstPort)

			if (tcp.DstPort == 1234 || tcp.DstPort == 5678 || tcp.DstPort == 9012) {
				if isKnockSequence(int(tcp.DstPort)) {
					fmt.Println("Correct knock sequence detected. Taking action...")
					takeAction()
					break
				}
			}
		}
	}
}

func isKnockSequence(port int) bool {

	fmt.Printf("Received knock on port: %d\n", port) // Debugging: Log the port of each received packet

	if port == knockSequence[currentKnockIndex] {
		fmt.Printf("Knock %d of the sequence detected on port %d\n", currentKnockIndex+1, port) // Debugging: Log correct sequence knock
		currentKnockIndex++
		if currentKnockIndex == len(knockSequence) {
			fmt.Println("Complete knock sequence detected.")
			// Reset for next time
			currentKnockIndex = 0
			return true
		}
	} else {
		fmt.Println("Incorrect knock sequence. Resetting.")
		// Reset if sequence is broken
		currentKnockIndex = 0
	}
	return false
}

func takeAction() {
	// Implement the action to be taken after the correct knock sequence
	fmt.Println("Action taken: Correct knock sequence detected.")
	// Here, you can define the action, such as opening a port, starting a service, etc.
}
