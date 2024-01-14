package main

import (
    "io/ioutil"
    "log"
    "path/filepath"
)

// SensitiveDirs holds the functionality to get sensitive directories
type SensitiveDirs struct{}

// GetUsers retrieves a list of user profiles from the system
func (sd *SensitiveDirs) GetUsers() ([]string, error) {
    var users []string
    entries, err := ioutil.ReadDir("C:\\Users")
    if err != nil {
        return nil, err
    }

    for _, entry := range entries {
        if entry.IsDir() {
            // Filter out default profiles and other non-user directories
            switch entry.Name() {
            case "Public", "Default", "Default User", "All Users", "desktop.ini":
                // Skip these directories
            default:
                users = append(users, entry.Name())
            }
        }
    }

    return users, nil
}

// Get returns a slice of sensitive directory paths, including user-specific paths
func (sd *SensitiveDirs) Get() []string {
    users, err := sd.GetUsers()
    if err != nil {
        log.Fatalf("Failed to get user profiles: %s", err)
    }

    dirs := []string{
        "C:\\Windows",
        "C:\\Windows\\System32",
        "C:\\Windows\\SysWOW64",
        "C:\\Windows\\Temp",
        "C:\\Windows\\Logs",
        "C:\\Windows\\Tasks",
        "C:\\Windows\\Prefetch",
        "C:\\Windows\\debug",
        "C:\\Windows\\tracing",
        "C:\\Program Files",
        "C:\\Program Files (x86)",
        "C:\\ProgramData",
        "C:\\Users",
        "C:\\Users\\Default",
        "C:\\Users\\Default User",
        "C:\\Users\\Public",
        "C:\\Documents and Settings",
        "C:\\Windows\\ServiceProfiles",
        "C:\\Windows\\ServiceProfiles\\LocalService",
        "C:\\Windows\\ServiceProfiles\\NetworkService",
        "C:\\Windows\\Panther",
        "C:\\Windows\\Inf",
        "C:\\Windows\\Fonts",
        "C:\\Windows\\Help",
        "C:\\Windows\\Media",
        "C:\\Windows\\Resources",
        "C:\\Windows\\Speech",
        "C:\\Windows\\Web",
        // Continue adding other general directories as needed
    }

    for _, user := range users {
        userDir := filepath.Join("C:\\Users", user)
        dirs = append(dirs, userDir)
        dirs = append(dirs, filepath.Join(userDir, "AppData"))
        dirs = append(dirs, filepath.Join(userDir, "AppData\\Local"))
        dirs = append(dirs, filepath.Join(userDir, "AppData\\Local\\Temp"))
        dirs = append(dirs, filepath.Join(userDir, "AppData\\Roaming"))
        dirs = append(dirs, filepath.Join(userDir, "Documents"))
        dirs = append(dirs, filepath.Join(userDir, "Desktop"))
        dirs = append(dirs, filepath.Join(userDir, "Downloads"))
        dirs = append(dirs, filepath.Join(userDir, "Pictures"))
        dirs = append(dirs, filepath.Join(userDir, "Music"))
        dirs = append(dirs, filepath.Join(userDir, "Videos"))
        dirs = append(dirs, filepath.Join(userDir, "Favorites"))
        dirs = append(dirs, filepath.Join(userDir, "Links"))
        dirs = append(dirs, filepath.Join(userDir, "Saved Games"))
        dirs = append(dirs, filepath.Join(userDir, "Contacts"))
        dirs = append(dirs, filepath.Join(userDir, "OneDrive"))
        dirs = append(dirs, filepath.Join(userDir, "Searches"))
        dirs = append(dirs, filepath.Join(userDir, "3D Objects"))
        dirs = append(dirs, filepath.Join(userDir, "Saved Pictures"))
        dirs = append(dirs, filepath.Join(userDir, "Saved Games"))
        dirs = append(dirs, filepath.Join(userDir, "Start Menu"))
        dirs = append(dirs, filepath.Join(userDir, "Templates"))
        // Continue adding other user-specific directories as needed
    }

    return dirs
}
