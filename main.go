package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"vault-tool/console"
	"vault-tool/dpapi"
)

var stdin = bufio.NewScanner(os.Stdin)

func readLine() (string, error) {
	if stdin.Scan() {
		return stdin.Text(), nil
	}
	if err := stdin.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func encryptDpapi(visual bool) {
	if visual {
		fmt.Println("Zadejte heslo: ")
	}
	secret, err := console.ReadPassword(readLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nepodařilo se načíst vstup: %v\n", err)
		os.Exit(1)
	}

	if len(secret) == 0 {
		fmt.Fprintln(os.Stderr, "Heslo nesmí být prázdné")
		os.Exit(1)
	}
	if visual {
		fmt.Printf("Uživatel zadal heslo o délce %d znaků\n", len(secret))
	}

	encrypted, err := dpapi.DpapiEncrypt([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "šifrování selhalo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("dpapi:" + base64.StdEncoding.EncodeToString(encrypted))
}

func decryptDpapi(visual bool) {
	if visual {
		fmt.Println("Zadejte zašifrovaná data - včetně dpapi: ")
	}

	input, err := readLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nepodařilo se načíst vstup: %v\n", err)
		os.Exit(1)
	}
	decrypted, err := dpapi.DpapiDecrypt(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Rozšifrování selhalo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(string(decrypted))
}

func decryptFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nelze načíst soubor %s: %v\n", path, err)
		os.Exit(1)
	}

	decrypted, err := dpapi.DpapiDecrypt(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nelze rozšifrovat soubor %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Println(decrypted)
}

func printHelp() {
	fmt.Println("vault-tool nástroj pro implementátory")
	fmt.Println("GORDIC, spol. s r.o. 2026")
	fmt.Println("v: 2.2")
	fmt.Println()
	fmt.Println("-help    nápověda")
	fmt.Println("-encrypt zašifruje heslo do dpapi:")
	fmt.Println("-decrypt rozšifruje dpapi:")
	fmt.Println("-decrypt [soubor] rozšifruje dpapi: ze souboru")
	fmt.Println("         některé windows terminály neumí vložit více jak 512 znaků!")
	fmt.Println()
	fmt.Println("bez parametrů se spustí interaktivní režim")
}

func executeInteractive() {
	fmt.Println("e\t-encrypt zašifruje")
	fmt.Println("d\t-decrypt rozšifruje")
	fmt.Println("q\t-konec")

	for {
		fmt.Println()
		fmt.Println("Zadejte příkaz: ")
		line, err := readLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Nepodařilo se načíst vstup: %v\n", err)
			os.Exit(1)
		}
		switch strings.TrimSpace(line) {
		case "d":
			decryptDpapi(true)
			return
		case "e":
			encryptDpapi(true)
			return
		case "q":
			return
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		executeInteractive()
		return
	}

	switch os.Args[1] {
	case "-help":
		printHelp()
		return
	case "-encrypt":
		encryptDpapi(false)
		return
	case "-decrypt":
		if len(os.Args) > 2 {
			decryptFile(os.Args[2])
		} else {
			decryptDpapi(false)
		}
		return
	}
}
