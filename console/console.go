package console

import (
	"fmt"
	"syscall"
	"unsafe"
)

func ReadPassword(readLine func() (string, error)) (string, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	handle, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return readLine()
	}
	var oldMode uint32
	r, _, _ := getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&oldMode)))
	if r == 0 {
		return readLine()
	}

	const enableEchoInput = 0x0004
	setConsoleMode.Call(uintptr(handle), uintptr(oldMode&^enableEchoInput))
	defer setConsoleMode.Call(uintptr(handle), uintptr(oldMode))

	line, err := readLine()
	fmt.Println()
	return line, err
}
