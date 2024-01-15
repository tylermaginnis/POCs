BITS 64

section .data
    title db 'Alert', 0
    message db 'Hello, World!', 0

section .text
    global main
    extern MessageBoxA, ExitProcess

main:
    sub rsp, 40                ; Shadow space for the Windows calling convention

    ; MessageBox(hwnd, text, caption, type)
    xor ecx, ecx               ; hWnd = NULL (no parent window)
    lea rdx, [message]         ; lpText
    lea r8, [title]            ; lpCaption
    mov r9d, 0                 ; uType = MB_OK
    call MessageBoxA

    ; Exit the process
    xor ecx, ecx               ; Return code 0
    call ExitProcess
