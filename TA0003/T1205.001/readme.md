# Network Knocking and Its Cybersecurity Significance

Network knocking, also known as port knocking, is a security practice with significant cybersecurity implications. This technique involves sending a sequence of connection attempts to predefined closed ports on a networked device, effectively serving as a digital secret handshake to access specific services. In this readme, we'll explore the concept of network knocking, its importance in cybersecurity, and how it can be employed both defensively and potentially exploitatively.

https://medium.com/@maginnist/proof-of-concept-a-deep-dive-into-port-knocking-7341ce9b8362

## Prerequisites

- **[npcap](https://npcap.com/)**
- **[gopacket](https://github.com/google/gopacket)**

## Understanding Network Knocking

- **Network Knocking Basics:** Network knocking conceals open ports by keeping them hidden until the correct knock sequence is received.
  
- **Concealing Open Ports:** The primary cybersecurity significance of network knocking is its ability to hide open ports on a device, reducing the attack surface.

## From a Defensive Perspective

- **Access Control:** Organizations can use network knocking to control and restrict access to critical services, allowing only users who know the knock sequence to gain access.

- **Reduced Attack Surface:** Concealing open ports mitigates the risk of unauthorized access, brute-force attacks, and port scanning.

- **Early Warning:** Failed or incorrect knock attempts can trigger alerts, providing early warning signs of potential threats.

## From an Adversarial Perspective

- **Concealing Backdoors:** Adversaries can use network knocking to hide backdoors or unauthorized access points on a compromised system. They may employ this technique to remain undetected and maintain control over a system.

- **Stealthy Reconnaissance:** Network knocking can enable attackers to discreetly scan a network for potential vulnerabilities without raising suspicion.

## Network Knock Scripts

These two scripts are designed for network knocking – a technique used to trigger specific actions on a network by sending a sequence of requests to predetermined ports. Network knocking enhances security by allowing access only to those who know the secret knock sequence.

### main.go

main.go is a Go program that captures incoming TCP packets on a selected network interface and checks if the packets match a predefined knock sequence. When the correct sequence is detected, it triggers a specified action.

#### Usage

1. Run the program and select a network interface.
2. Wait for the knock sequence (default: 1234, 5678, 9012) on the selected interface.
3. When the correct sequence is detected, it will take the specified action (by default, it prints "Correct knock sequence detected").

```bash
go mod tidy
go run main.go
```

## How it Works

### main.go

The `main.go` program uses the `pcap` library to capture packets. It functions as follows:

1. It listens for incoming TCP packets.
2. It checks if the destination port matches the predefined sequence.
3. When a valid knock sequence is detected, it triggers the `takeAction` function.

### knocker.go

`knocker.go` is a Go script designed to send TCP connection attempts to specified IP addresses and ports. It simulates the act of "knocking" on ports.

#### Usage

1. Run the program.
2. Enter the IP address to knock on.
3. Enter three port numbers to knock on.
4. The program will attempt to establish a TCP connection to the specified ports on the provided IP address.

```go run knocker.go
```

## How it Works

The script operates as follows:

- It utilizes the `net` package to establish TCP connections.
- The script attempts to connect to the specified IP address and ports.
- If a connection is successful, it signifies that the port is "reachable."

These scripts are made available for educational purposes and can be customized to execute actions beyond simple print statements. Always exercise responsible use of network-knocking techniques and operate within legal boundaries.

## Conclusion

Network knocking is a cybersecurity technique with the potential to significantly enhance network security when deployed correctly. However, it should be approached with caution, and robust security measures must be in place to prevent misuse. Network knocking serves as an additional layer of defense by concealing open ports, managing access, and functioning as an early warning system. Nevertheless, the dual nature of network knocking underscores the importance of a balanced approach to network security.

Please exercise diligence and carefully consider the implications of network knocking in your specific use case. Utilize it wisely to reinforce your network security.

*Disclaimer: The information provided here is intended for educational purposes, and any implementation should adhere to legal and ethical standards.*
