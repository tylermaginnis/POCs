# Proof of Concepts (POCs) Index

This document serves as an index for the various Proof of Concepts (POCs) developed as part of my cybersecurity research and testing.

## List of POCs

### Proof of Concept: NTUSER.DAT Registry Extractor

This Go Proof of Concept (POC) demonstrates how to extract information from Windows NTUSER.DAT files using Go and PowerShell commands. NTUSER.DAT files are part of Windows user profiles and contain user-specific registry settings.

- [TA0007/T087.001](./TA0007/T087.001): Account Discovery Technique
  - This script is designed to demonstrate the Account Discovery technique. [Read more](./TA0007/T087.001/readme.md).

### Proof of Concept: Network Information Gathering Tool

This Go program gathers network information from your system and displays it in JSON format. It provides details about network interfaces and ARP cache entries.

- [TA0102/T0846](./TA0102/T0846): Network Information Gathering Tool
  - This tool retrieves network interface and ARP cache information. [Read more](./TA0102/T0846/readme.md).

### Proof of Concept: Automated Collection

This Go program monitors and logs keyboard and mouse input events on a Windows system. It captures keypresses and mouse clicks and saves them to an output file. Additionally, it can take screenshots when a mouse click event occurs.

- [TA0009/T1119](./TA0009/T1119/readme.md): Automated Collection
  - This tool demonstrates input monitoring techniques. [Read more](./TA0009/T1119/readme.md).

## Contributing

Please follow the [contribution guidelines](./CONTRIBUTING.md) to propose changes or add new POCs.

## Disclaimer

All POCs in this repository are for educational and authorized testing purposes only. Ensure you have the necessary permissions before running any of these scripts.

## License

MIT 3.0

---

For any queries or contributions, please contact via GH.
