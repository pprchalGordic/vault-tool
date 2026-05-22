package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("vault-tool nástroj pro implementátory")
		fmt.Println("GORDIC, spol. s r.o. 2026")
		fmt.Println("v: 1.4")
		fmt.Println()
		fmt.Println("-encrypt zašifruje heslo do dpapi: tvaru")
		fmt.Println("-decrypt rozšifruje dpapi: hodnotu")
		return
	}

	switch os.Args[1] {
	case "-encrypt":
		encryptInteractive()
		return
	case "-decrypt":
		path := ""
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		decryptInteractive(path)
		return
	}
}
