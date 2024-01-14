# Input Logger

## Overview

The Input Logger is a Go program designed to monitor and log keyboard and mouse input events on a Windows system. It captures keypresses and mouse clicks and saves them to an output file. Additionally, it can take screenshots when a mouse click event occurs.

https://medium.com/@maginnist/unlocking-windows-automated-collection-a-deep-dive-into-the-automated-collection-poc-934587555b21

## Prerequisites

Before using this tool, ensure that you have the following prerequisites:

- A Windows operating system
- Administrative privileges to execute the program
- Go programming language installed

## Usage

1. Clone or download the project.

2. Open a terminal or command prompt.

3. Navigate to the project directory.

4. Run the Go program using the following command:

   ```bash
   go run main.go
   ```

## Features

### Monitor Keyboard Input

- Captures keypress events.
- Logs the virtual key code of the pressed key.

### Monitor Mouse Input

- Captures mouse click events (left and right buttons).
- Logs the coordinates (X, Y) where the click occurred.

### Capture Screenshots

- Takes a screenshot when a mouse click event is detected.
- Saves the screenshot as a PNG file with a timestamp.

## Notes

- This tool is for educational and informational purposes only.
- It is designed to demonstrate input monitoring techniques.
- Ensure you have administrative privileges to execute the program.
- The program currently logs input and takes screenshots to an output file named "input_log.txt."
- For more advanced functionality, you can implement actual CallNextHookEx from user32.dll.

## License

This project is licensed under the MIT License. Feel free to modify and use it as needed.

**Please note:** The input logger may capture sensitive information, so use it responsibly and only on systems you have permission to access.

