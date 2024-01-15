package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the IP address to knock on: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	ports := make([]string, 3)
	for i := 0; i < 1; i++ {
		fmt.Printf("#%d: Enter port number to knock on: ", i+1)
		port, _ := reader.ReadString('\n')
		ports[i] = strings.TrimSpace(port)
	}

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		fmt.Printf("Knocking on %s with custom packet...\n", address)

		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			fmt.Printf("Failed to knock on %s: %s\n", address, err)
			continue
		}

		// Craft the custom packet with standard HTTP headers and the special string
		customPacket := "GET / HTTP/1.1\r\nHost: " + ip + "\r\nUser-Agent: CustomKnockClient\r\nCustom-Data: custom_packet\r\n\r\n"
		_, err = conn.Write([]byte(customPacket))
		if err != nil {
			fmt.Printf("Failed to send custom packet to %s: %s\n", address, err)
		} else {
			fmt.Printf("Custom packet sent to %s successfully.\n", address)
		}

		conn.Close()
	}
}
