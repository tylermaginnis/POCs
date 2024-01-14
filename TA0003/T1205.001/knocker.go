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
    for i := 0; i < 3; i++ {
        fmt.Printf("Enter port number %d to knock on: ", i+1)
        port, _ := reader.ReadString('\n')
        ports[i] = strings.TrimSpace(port)
    }

    for _, port := range ports {
        address := fmt.Sprintf("%s:%s", ip, port)
        fmt.Printf("Knocking on %s...\n", address)

        conn, err := net.DialTimeout("tcp", address, 5*time.Second)
        if err != nil {
            if strings.Contains(err.Error(), "refused") {
                fmt.Printf("Knock on %s sent (connection refused as expected).\n", address)
            } else {
                fmt.Printf("Failed to knock on %s: %s\n", address, err)
            }
        } else {
            conn.Close()
            fmt.Printf("Knocked on %s successfully.\n", address)
        }
    }
}