package main

import (
    "crypto/rand"
    "fmt"
    "log"
    "math/big"
    "os"
    "path/filepath"
    "time"
    mrand "math/rand"
)

// GenerateRandomBytes returns securely generated random bytes of given length
func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }

    return b, nil
}

// CreateRandomFile creates a file with random data and a given extension
func CreateRandomFile(dir, ext string, size int) (string, error) {
    data, err := GenerateRandomBytes(size)
    if err != nil {
        return "", err
    }

    // Generate a random integer for the file name
    randInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
    if err != nil {
        return "", err
    }

    filePath := filepath.Join(dir, fmt.Sprintf("suspicious_%s%s", randInt.String(), ext))
    file, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    _, err = file.Write(data)
    return filePath, err
}

func main() {
    // Seed the random number generator
    mrand.Seed(time.Now().UnixNano())

    // Initialize SensitiveDirs and get the list of directories
    sd := SensitiveDirs{}
    dirs := sd.Get()

    // List of extensions
    extensions := []string{".bin", ".exe", ".dll"}

    // Size of the random file in bytes
    fileSize := 1024 // 1KB

    // Creating output file
    timestamp := time.Now().Format("20060102-150405")
    reportFile, err := os.Create(fmt.Sprintf("sensitiveFileCreationReport_%s.txt", timestamp))
    if err != nil {
        log.Fatalf("Failed to create report file: %s", err)
    }
    defer reportFile.Close()

    var createdFiles []string

    for _, dir := range dirs {
        for _, ext := range extensions {
            filePath, err := CreateRandomFile(dir, ext, fileSize)
            if err != nil {
                fmt.Fprintf(reportFile, "Failed to create file with extension %s in %s: %s\n", ext, dir, err)
            } else {
                fmt.Fprintf(reportFile, "File with extension %s created successfully in %s\n", ext, dir)
                createdFiles = append(createdFiles, filePath)
            }
        }
    }

    // Clean up created files
    for _, file := range createdFiles {
        err := os.Remove(file)
        if err != nil {
            fmt.Fprintf(reportFile, "Failed to delete file: %s, error: %s\n", file, err)
        } else {
            fmt.Fprintf(reportFile, "Deleted file: %s\n", file)
        }
    }
}
