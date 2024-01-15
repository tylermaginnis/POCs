# Socket Filtering Proof of Concept

This repository contains two Go scripts, `main.go` and `knocker.go`, that demonstrate socket filtering techniques. Socket filtering involves capturing and processing network packets to filter and act upon specific types of traffic. The scripts showcase both the defensive and adversarial aspects of socket filtering.

## Prerequisites

Before running the scripts, make sure you have the following prerequisites installed:

- [gopacket](https://github.com/google/gopacket)
- [npcap](https://npcap.org)

## main.go

`main.go` is a Go program that captures incoming TCP packets on a selected network interface and checks for custom packets on a specific target port (default: port 80). When a custom packet is detected, it triggers a specified action.

### Usage

1. Run the program and select a network interface.
2. Wait for custom packets on the target port (default: 80) on the selected interface.
3. When a custom packet is detected, it takes the specified action (by default, it prints "Custom packet detected.").

```shell
go mod tidy
go run main.go
```

## How it Works

### main.go

The `main.go` program utilizes the `gopacket` library to capture network packets. Here's how it works:

1. It listens for incoming TCP packets on a selected network interface.
2. The script checks if the destination port matches the target port (default: port 80).
3. When a valid custom packet is detected, it triggers the `takeAction` function.

### knocker.go

`knocker.go` is a Go script that simulates port knocking. Here's how it works:

1. Run the program.
2. Enter the IP address to knock on.
3. Enter the port number to knock on.
4. The script attempts to establish a TCP connection to the specified IP address and port with a custom packet.

**Usage**

To use `knocker.go`:

1. Run the program.
2. Enter the IP address to target.
3. Enter the port number to knock on.
4. The program attempts to establish a TCP connection to the specified IP address and port using a custom packet.

These scripts provide practical demonstrations of socket filtering techniques. Modify and utilize them responsibly for your specific use cases.

```shell
go run knocker.go
```

## How it Works

The provided scripts demonstrate socket filtering techniques as follows:

### main.go

1. The `main.go` script uses the `net` package to create TCP connections.
2. It attempts to connect to the specified IP address and port.
3. If a connection is successful, it sends a custom packet containing specific headers and a special string.

These scripts are intended for educational purposes and can be customized to perform actions beyond simple print statements. Always exercise caution and ensure compliance with legal and ethical standards when employing socket filtering techniques.

## Conclusion

Socket filtering techniques play a critical role in network security. They can be employed defensively, as demonstrated in `main.go`, to monitor and filter incoming traffic. Additionally, adversaries can utilize similar techniques to send custom packets, as showcased in `knocker.go`. Understanding both defensive and adversarial aspects of socket filtering is essential for enhancing network security.

Please exercise diligence and carefully consider the implications of socket filtering in your specific use case. Use these techniques judiciously to bolster your network security.

Disclaimer: The information provided here is for educational purposes, and any practical implementation should adhere to legal and ethical standards.
