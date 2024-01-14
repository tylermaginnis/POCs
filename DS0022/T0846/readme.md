# Network Information Gathering Tool

## Overview

This Go program is designed to gather network information from your system and display it in JSON format. It provides details about network interfaces and ARP (Address Resolution Protocol) cache entries.

## Prerequisites

Before using this tool, ensure that you have the following prerequisites:

- [Go programming language](https://golang.org/dl/)
- Administrative privileges to execute network-related commands

## Usage

1. Clone or download the project.

2. Open a terminal or command prompt.

3. Navigate to the project directory.

4. Run the Go program using the following command:

   ```bash
   go run main.go
   ```

## Features

### Gather Network Interface Information

- Retrieves information about all network interfaces available on your system.
- For each network interface, collects the following details:
  - Interface Name
  - IP Addresses
  - Hostnames associated with the IP addresses (if available)

### Gather ARP Cache Information

- Retrieves information about the ARP cache on your system.
- For each ARP cache entry, collects the following details:
  - IP Address
  - MAC Address
  - Hostnames associated with the IP address (if available)

### Display JSON Output

- Combines the gathered information into a structured JSON output.
- Displays the JSON output in the terminal.

## Notes

- This tool is designed for informational and diagnostic purposes.
- Ensure you have administrative privileges to execute network-related commands.
- Network information may vary depending on your system's configuration.
- Hostnames associated with IP addresses are resolved if possible.

## License

This project is licensed under the MIT License. Feel free to modify and use it as needed.

**Please note:** Network information may contain sensitive data, so use this tool responsibly and only on systems you have permission to access.
