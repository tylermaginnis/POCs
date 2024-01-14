# UDP - LLMNR & NBT-NS Response Simulator

The UDP Response Simulator is a Go program designed to simulate responses for LLMNR (Link-Local Multicast Name Resolution) and NBT-NS (NetBIOS Name Service) protocols. It listens on specific UDP ports and sends simulated responses to incoming requests.

## Usage

1. Clone or download the project.

2. Open a terminal or command prompt.

3. Navigate to the project directory.

4. Run the Go program using the following command:

   ```bash
   go run main.go
   ```

# UDP Response Simulator

## Overview

The program listens on two UDP ports:

- **LLMNR (Port 5355):** Simulates LLMNR responses.
- **NBT-NS (Port 137):** Simulates NBT-NS responses.

It constructs simulated response packets based on the respective protocol specifications and sends them to the requesting clients.

## Simulated Responses

### Simulated LLMNR Response

The program constructs a simulated LLMNR response packet. You can customize the hostname and IPv6 address in the response:

- **Hostname:** "example.com" (You can replace this with your desired hostname)
- **IPv6 Address:** 2001:0db8:0000:0000:0000:0000:0000:0001 (Replace with the desired IPv6 address)

### Simulated NBT-NS Response

The program constructs a simulated NBT-NS response packet. You can customize the NetBIOS name and IPv4 address in the response:

- **NetBIOS Name:** "MYCOMPUTER" (Replace with your desired NetBIOS name)
- **IPv4 Address:** 192.168.1.100 (Replace with the desired IPv4 address)

Please note that these are simplified and simulated response packets for demonstration purposes. In a real-world scenario, you would need to construct response packets based on the protocol specifications and requirements.

## License

This project is licensed under the MIT License. Feel free to modify and use it as needed.

**Please note:** This program is intended for educational and demonstration purposes. Use it responsibly and only on systems you have permission to access.
