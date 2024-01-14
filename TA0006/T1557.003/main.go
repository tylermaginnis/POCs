package main

import (
    "fmt"
    "net"
    "sync"
)

var (
    ipPool      []net.IP
    ipPoolMutex sync.Mutex
)

func main() {
    // Initialize the IP pool
    initializeIPPool("192.168.1.100", "192.168.1.150")

    // Start listening for DHCP DISCOVER packets
    dhcpListener()
}

func initializeIPPool(startIP, endIP string) {
    ipPoolMutex.Lock()
    defer ipPoolMutex.Unlock()

    start := net.ParseIP(startIP).To4()
    end := net.ParseIP(endIP).To4()
    for ip := start; !ip.Equal(end); ip = incrementIP(ip) {
        ipPool = append(ipPool, ip)
    }
    ipPool = append(ipPool, end) // Include the end IP
}

func dhcpListener() {
    conn, err := net.ListenPacket("udp4", ":67")
    if err != nil {
        fmt.Println("Error setting up UDP listener:", err)
        return
    } else {
        fmt.Printf("Setup UDP Listener on Port 67")
    }
    defer conn.Close()

    for {
        buffer := make([]byte, 1500) // Standard MTU size
        _, addr, err := conn.ReadFrom(buffer)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            continue
        }

        if isDHCPDiscover(buffer) {
            fmt.Printf("DHCP Discover received from %v\n", addr)
            response := craftDHCPResponse(buffer)
            if response != nil {
                conn.WriteTo(response, addr)
                fmt.Printf("Wrote malicious response")
            }
        }
    }
}

func isDHCPDiscover(packet []byte) bool {
    // Ensure packet length is sufficient for DHCP (minimum 240 bytes + options)
    if len(packet) < 244 {
        return false
    }

    // Check if the message type is BOOTREQUEST (1)
    if packet[0] != 1 {
        return false
    }

    // DHCP options start at byte 240
    options := packet[240:]

    for i := 0; i < len(options); {
        optionType := options[i]
        i++

        // End Option
        if optionType == 255 {
            break
        }

        optionLen := int(options[i])
        i++

        // DHCP Message Type option
        if optionType == 53 && optionLen == 1 && options[i] == 1 { // 1 = DHCP Discover
            return true
        }

        i += optionLen
    }

    return false
}



func craftDHCPResponse(discoverPacket []byte) []byte {
    // For demonstration, let's assign static IPs
    offeredIP := net.ParseIP("192.168.1.200").To4()
    serverIP := net.ParseIP("192.168.1.1").To4()

    response := make([]byte, 300)
    copy(response, discoverPacket)

    // Set message type to BOOTREPLY
    response[0] = 2

    // Set your server IP
    copy(response[20:24], serverIP)

    // Set the offered IP
    copy(response[16:20], offeredIP)

    // DHCP options
    response[240] = 53 // DHCP Message Type
    response[241] = 1  // Length
    response[242] = 2  // 2 = DHCP Offer
    response[243] = 255 // End Option

    return response
}



func incrementIP(ip net.IP) net.IP {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
    return ip
}
