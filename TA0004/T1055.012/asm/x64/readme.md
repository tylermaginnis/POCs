# TA0002\T1047\asm\x64\payload.asm

This file contains a HelloWorld call via MessageBoxA and ExitProcess through W32 API.

## Prerequisites

- [NASM] (https://www.nasm.us/)
- [link.exe] (https://learn.microsoft.com/en-us/cpp/build/reference/linking?view=msvc-170)

### Compilation

```bash
	nasm -f win64 payload.asm -o payload.obj
```

### Linking

To produce an EXE:

```bash
	link /subsystem:console /entry:main /LARGEADDRESSAWARE:NO payload.obj User32.lib Kernel32.lib
```
