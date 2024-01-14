package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func main() {
    userProfilesDir := "C:\\Users"

    entries, err := os.ReadDir(userProfilesDir)
    if err != nil {
        fmt.Println("Error reading directory:", err)
        return
    }

    for _, entry := range entries {
        if entry.IsDir() {
            userName := entry.Name()
            userDir := filepath.Join(userProfilesDir, userName)
            ntuserDatPath := filepath.Join(userDir, "NTUSER.DAT")

            if _, err := os.Stat(ntuserDatPath); err == nil {
                fmt.Printf("Found NTUSER.DAT for user '%s'\n", userName)
                handleNTUSERDat(userName, ntuserDatPath)
            }
        }
    }
}

func handleNTUSERDat(userName, ntuserDatPath string) {
    regKeyPath := "HKLM\\TempHive_" + userName

    loadCmd := fmt.Sprintf("reg.exe load %s %s", regKeyPath, ntuserDatPath)
    if err := executePowerShellCommand(loadCmd); err != nil {
        fmt.Printf("Failed to load NTUSER.DAT for user '%s': %v\n", userName, err)
        return
    }

    // Recursively read the loaded hive
    readCmd := fmt.Sprintf("Get-ChildItem -Path 'Registry::%s' -Recurse", regKeyPath)
    if err := executePowerShellCommand(readCmd); err != nil {
        fmt.Printf("Failed to read NTUSER.DAT for user '%s': %v\n", userName, err)
    }

    unloadCmd := fmt.Sprintf("reg.exe unload %s", regKeyPath)
    if err := executePowerShellCommand(unloadCmd); err != nil {
        fmt.Printf("Failed to unload NTUSER.DAT for user '%s': %v\n", userName, err)
        return
    }
}

func executePowerShellCommand(command string) error {
    cmd := exec.Command("powershell", "-Command", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("error executing command: %s, output: %s, error: %w", command, output, err)
    }

    fmt.Println(string(output))
    return nil
}
