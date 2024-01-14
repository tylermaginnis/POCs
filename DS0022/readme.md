# Sensitive File Creation and Cleanup Script

## Overview

This script is designed for cybersecurity testing purposes. It creates random files with specified extensions in a set of sensitive Windows directories and then cleans them up. The script generates a report detailing the creation and deletion of these files.

## Features

- **File Creation**: Creates random data files with extensions `.bin`, `.exe`, and `.dll` in various sensitive directories.
- **Report Generation**: Outputs a report named `sensitiveFileCreationReport_TIMESTAMP.txt`, logging the file creation and deletion process.
- **Cleanup**: Removes all created files after the process.

## Components

The script consists of two main files:

- `main.go`: The main script that handles file creation, logging, and cleanup.
- `sensitiveDirs.go`: A supporting file that provides a list of sensitive directories.

## Usage

1. **Prerequisites**:
   - Ensure you have Go installed on your system.
   - This script should only be run on systems where you have explicit permission for such testing.

2. **Running the Script**:
   - Place `main.go` and `sensitiveDirs.go` in the same directory.
   - Run the script using the command: `go run main.go sensitiveDirs.go`.

## Output

- The script will create a report file in the format `sensitiveFileCreationReport_YYYYMMDD-HHMMSS.txt`.
- The report logs the success or failure of file creation and deletion in each directory.

## Disclaimer

This script is intended for educational and authorized penetration testing purposes only. It should be used responsibly and ethically. Running this script on systems without proper authorization may violate privacy and security policies and could be illegal.

## Author

@tylermaginnis

## License

MIT 3.0

---

**Note**: The directories and actions performed by this script are potentially intrusive. It's crucial to understand the system impact and have necessary backups and permissions.
