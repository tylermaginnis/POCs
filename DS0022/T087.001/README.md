# Go POC: Extract Information from Windows NTUSER.DAT

This Go Proof of Concept (POC) demonstrates how to extract information from Windows NTUSER.DAT files using Go and PowerShell commands. NTUSER.DAT files are part of Windows user profiles and contain user-specific registry settings.

## Prerequisites

Before running the POC, ensure that you have the following prerequisites installed:

- Go programming language: [https://golang.org/dl/](https://golang.org/dl/)
- PowerShell (comes pre-installed on most Windows systems)
- Administrative privileges to execute `reg.exe` commands

## Usage

1. Clone or download the project.

2. Open a terminal or command prompt.

3. Navigate to the project directory.

4. Run the Go program using the following command:

```bash
go run main.go
```

This program will perform the following actions:

1. Scan the user profiles directory (C:\Users) and look for NTUSER.DAT files in each user's directory.

2. For each NTUSER.DAT file found, it will:
    - Load the NTUSER.DAT file into the Windows registry under a temporary registry key.
    - Recursively read the loaded registry hive.
    - Unload the registry hive after reading.

3. The program will display the registry information it reads from each NTUSER.DAT file.

4. The program will fail to display the registry information of the current user, due to file locking.

## Notes

- This POC is designed for educational purposes and may need modification for production use.
- Ensure you have administrative privileges to execute `reg.exe` commands.
- Be cautious when working with the Windows registry, as it contains critical system settings.

## License

This project is licensed under the MIT License.
