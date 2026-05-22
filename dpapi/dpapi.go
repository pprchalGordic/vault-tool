package dpapi

import (
	"encoding/base64"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

const cryptProtectLocalMachine = 0x4

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func DpapiEncrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("prázdná data")
	}
	input := dataBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}
	var output dataBlob

	r, _, err := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(&input)),
		0, 0, 0, 0, cryptProtectLocalMachine,
		uintptr(unsafe.Pointer(&output)),
	)
	if r == 0 {
		return nil, fmt.Errorf("Nepodařilo se zašifrovat data: %v", err)
	}

	result := make([]byte, output.cbData)
	copy(result, unsafe.Slice(output.pbData, output.cbData))
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(output.pbData)))
	return result, nil
}

func DpapiDecrypt(input string) (string, error) {
	encoded := strings.TrimPrefix(strings.TrimSpace(input), "dpapi:")
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("Není ve formátu base64: %v", err)
	}

	decrypted, err := dpapiDecrypt(data)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func dpapiDecrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("prázdná data")
	}
	input := dataBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}
	var output dataBlob

	r, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&input)),
		0, 0, 0, 0, cryptProtectLocalMachine,
		uintptr(unsafe.Pointer(&output)),
	)
	if r == 0 {
		return nil, fmt.Errorf("Nepodařilo se rozšifrovat data: %v", err)
	}

	result := make([]byte, output.cbData)
	copy(result, unsafe.Slice(output.pbData, output.cbData))
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(output.pbData)))
	return result, nil
}
