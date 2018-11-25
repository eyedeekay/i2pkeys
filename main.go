package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/eyedeekay/sam3"
)

var (
	keyfile        = flag.String("keyfile", "default.i2pkeys", "key file to use for both generation and reading")
	shortkeyfile   = flag.String("file", "default.i2pkeys", "short for 'keyfile'")
	shorterkeyfile = flag.String("f", "default.i2pkeys", "shorter for 'file'")
	usestdin       = flag.Bool("input", false, "read key base64 from stdin")
	useshortstdin  = flag.Bool("i", false, "short for 'input'")
	generate       = flag.Bool("generate", true, "generate new keys(requires SAM connection for now)")
	generatekeys   = flag.Bool("g", true, "short for 'generate'")
	samaddr        = flag.String("samaddress", "127.0.0.1:7656", "")
	sam            = flag.String("s", "127.0.0.1:7656", "short for 'samaddress'")
	delimiter      = flag.String("delimiter", "=", "string to use as a delimiter in output")
	delim          = flag.String("d", "=", "short for 'delimiter'")
)

func main() {
	flag.Parse()
	if *delim != "=" || *delimiter != "=" {
		if *delimiter != "=" {
			delimiter = delimiter
		} else if *delim != "=" {
			delimiter = delim
		}
	}
	if *usestdin || *useshortstdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			str := scanner.Text()
			addr, err := sam3.NewI2PAddrFromString(str)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("base32%s%s", *delimiter, addr.Base64())
			fmt.Printf("base64%s%s", *delimiter, addr.Base32())
		}
		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}
	} else if *keyfile != "" || *shortkeyfile != "" || *shorterkeyfile != "" {
		if *keyfile != "" {
			keyfile = keyfile
		} else if *shortkeyfile != "" {
			keyfile = shortkeyfile
		} else if *shorterkeyfile != "" {
			keyfile = shorterkeyfile
		}
		var dogeneratekeys = false
		if *generate || *generatekeys {
			dogeneratekeys = true
		}
		if dogeneratekeys {
			GenerateKeys()
		} else {
			DisplayKeys()
		}
	}
}

func GenerateKeys() {
	var samaddrstring = "127.0.0.1:7657"
	if *samaddr != "" {
		samaddrstring = *samaddr
	} else if *sam != "" {
		samaddrstring = *sam
	}
	sam, err := sam3.NewSAM(samaddrstring)
	if err != nil {
		panic(err.Error())
	}
	keys, err := sam.NewKeys()
	if err != nil {
		panic(err.Error())
	}

	openfile, err := os.Open(*keyfile)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			openfile, err = os.Create(*keyfile)
		} else {
			panic(err.Error())
		}
	}
	if err != nil {
		panic(err.Error())
	}
	err = sam3.StoreKeysIncompat(keys, openfile)
	if err != nil {
		panic(err.Error())
	}
	openfile.Close()
}

func DisplayKeys() {
	openfile, err := os.Open(*keyfile)
	if err != nil {
		panic(err.Error())
	}
	addr, err := sam3.LoadKeysIncompat(openfile)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("base32%s%s", *delimiter, addr.Addr().Base64())
	fmt.Printf("base64%s%s", *delimiter, addr.Addr().Base32())
	openfile.Close()
}
