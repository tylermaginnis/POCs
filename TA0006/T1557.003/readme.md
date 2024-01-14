# DHCP Spoofing Proof of Concept

This Go program is a proof of concept for DHCP spoofing. It demonstrates how to listen for DHCP DISCOVER packets on UDP port 67 and respond with a malicious DHCP OFFER packet.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Usage](#usage)
- [Understanding the Code](#understanding-the-code)
- [Legal and Ethical Considerations](#legal-and-ethical-considerations)

## Prerequisites

Before running this program, make sure you have the following:

- Go programming environment set up.
- Appropriate permissions to run network-related code.
- A controlled and isolated network environment for testing.

## Usage

1. Clone or download this repository to your local machine.

2. Navigate to the project directory using the command line.

3. Run the following command to execute the program:

   ```bash
   go run main.go
   ```

## Understanding the Code

The program includes the following key functions:

- `initializeIPPool(startIP, endIP string)`: Initializes an IP pool to assign IP addresses to DHCP clients.
- `dhcpListener()`: Listens for DHCP DISCOVER packets and sends malicious DHCP OFFER responses.
- `isDHCPDiscover(packet []byte) bool`: Checks if the received packet is a DHCP DISCOVER packet.
- `craftDHCPResponse(discoverPacket []byte) []byte`: Crafts a malicious DHCP OFFER response packet.
- `incrementIP(ip net.IP) net.IP`: Increments an IP address by 1.

## Legal and Ethical Considerations

- This program is intended for educational and research purposes only.
- Use it only in a controlled and authorized network environment.
- Be aware of the potential legal and ethical implications of DHCP spoofing.
- Always obtain proper authorization before conducting any network-related tests.

## License

This project is licensed under the MIT License. Feel free to modify and use it as needed

**Note:** The use of DHCP spoofing outside of controlled environments and without proper authorization can be illegal and unethical. Always act responsibly and with respect to network security and privacy.
