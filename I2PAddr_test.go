package i2pkeys

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	yoursam               = "127.0.0.1:7656"
	validShortenedI2PAddr = "idk.i2p"
	validI2PAddrB32       = "b2o47zwxqjbn7jj37yqkmvbmci7kqubwgxu3umqid7cexmc7xudq.b32.i2p"
	validI2PAddrB64       = "spHxea2xhPjKH9yyEeFJ96aqtvKidH-GiWxs8dH6RWS2FrDoWFhuEkfw77pF~Hv57lLhMaMB3qqWjCtYXOjL48Q1zYbr3MAcTO44wwVPjOU1hU77vbJcUuwBeRvaSr2dZx-FiTSOdQuhPD1EozYNRIMFwZ0fZwKf~3Gj4dEWccOLKs~NbiPsj-~tc5tmhAs8yBeoZEqEBe40X75SfSHY-EnstcZevVAwIXYk3zX3KF0mji3bo2QXuTFcMZHHLiLd2AHLRANzWyvQ9DC1rnCsHJM4xxV4dVp0pHkP1hwBo7E0NJvN4nFkQcj-FI2RJ~cFUCk7qc86PRHwvKCjzSlrgjtDsMUwd83Dz1PfpzCqHNLUFWI7uPKbKcJZhasFm4kEhUyupd85q75Ch2IZE9J2JXodSxmseO5ZKcHK6pFtfR-HbzKjIe92TWHsNkmvtoHiUaOVrWnk-cmo2I1W1VxfL08teDxQ13P80uFaMcameRzuFM2F8pSOpoyEJUDRGLEeBQAEAAcAAA=="
)

func Test_Basic(t *testing.T) {
	fmt.Println("Test_Basic")
	fmt.Println("\tAttaching to SAM at " + yoursam)
	keys, err := NewDestination()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(keys.String())
}

func Test_Basic_Lookup(t *testing.T) {
	fmt.Println("Test_Basic")
	fmt.Println("\tAttaching to SAM at " + yoursam)
	keys, err := Lookup(validShortenedI2PAddr)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(keys.String())
}

func Test_NewI2PAddrFromString(t *testing.T) {
	t.Run("Valid base64 address", func(t *testing.T) {
		addr, err := NewI2PAddrFromString(validI2PAddrB64)
		if err != nil {
			t.Fatalf("NewI2PAddrFromString failed for valid address: '%v'", err)
		}
		if addr.Base64() != validI2PAddrB64 {
			t.Errorf("NewI2PAddrFromString returned incorrect address. Got '%s', want '%s'", addr.Base64(), validI2PAddrB64)
		}
	})

	t.Run("Invalid address", func(t *testing.T) {
		invalidAddr := "not-a-valid-address"
		_, err := NewI2PAddrFromString(invalidAddr)
		if err == nil {
			t.Error("NewI2PAddrFromString should have failed for invalid address")
		}
	})

	t.Run("Base32 address", func(t *testing.T) {
		_, err := NewI2PAddrFromString(validI2PAddrB32)
		if err == nil {
			t.Error("NewI2PAddrFromString should have failed for base32 address")
		}
	})

	t.Run("Empty address", func(t *testing.T) {
		_, err := NewI2PAddrFromString("")
		if err == nil {
			t.Error("NewI2PAddrFromString should have failed for empty address")
		}
	})

	t.Run("Address with .i2p suffix", func(t *testing.T) { //CHECK
		addr, err := NewI2PAddrFromString(validI2PAddrB64 + ".i2p")
		if err != nil {
			t.Fatalf("NewI2PAddrFromString failed for address with .i2p suffix: '%v'", err)
		}
		if addr.Base64() != validI2PAddrB64 {
			t.Errorf("NewI2PAddrFromString returned incorrect address. Got '%s', want '%s'", addr.Base64(), validI2PAddrB64)
		}
	})
}

func Test_I2PAddr(t *testing.T) {
	addr := I2PAddr(validI2PAddrB64)
	base32 := addr.Base32()

	t.Run("Base32 suffix", func(t *testing.T) {
		if !strings.HasSuffix(base32, ".b32.i2p") {
			t.Errorf("Base32 address should end with .b32.i2p, got %s", base32)
		}
	})

	t.Run("Base32 length", func(t *testing.T) {
		if len(base32) != 60 {
			t.Errorf("Base32 address should be 60 characters long, got %d", len(base32))
		}
	})
}

func Test_DestHashFromString(t *testing.T) {
	t.Run("Valid hash", func(t *testing.T) {
		hash, err := DestHashFromString(validI2PAddrB32)
		if err != nil {
			t.Fatalf("DestHashFromString failed for valid hash: '%v'", err)
		}
		if hash.String() != validI2PAddrB32 {
			t.Errorf("DestHashFromString returned incorrect hash. Got '%s', want '%s'", hash.String(), validI2PAddrB32)
		}
	})

	t.Run("Invalid hash", func(t *testing.T) {
		invalidHash := "not-a-valid-hash"
		_, err := DestHashFromString(invalidHash)
		if err == nil {
			t.Error("DestHashFromString should have failed for invalid hash")
		}
	})

	t.Run("Empty hash", func(t *testing.T) {
		_, err := DestHashFromString("")
		if err == nil {
			t.Error("DestHashFromString should have failed for empty hash")
		}
	})
}

func Test_I2PAddrToBytes(t *testing.T) {
	addr := I2PAddr(validI2PAddrB64)

	t.Run("ToBytes and back", func(t *testing.T) {
		decodedBytes, err := addr.ToBytes()
		if err != nil {
			t.Fatalf("ToBytes failed: '%v'", err)
		}

		encodedString := i2pB64enc.EncodeToString(decodedBytes)
		if encodedString != validI2PAddrB64 {
			t.Errorf("Round-trip encoding/decoding failed. Got '%s', want '%s'", encodedString, validI2PAddrB64)
		}
	})

	t.Run("Direct decoding comparison", func(t *testing.T) {
		decodedBytes, err := addr.ToBytes()
		if err != nil {
			t.Fatalf("ToBytes failed: '%v'", err)
		}

		directlyDecoded, err := i2pB64enc.DecodeString(validI2PAddrB64)
		if err != nil {
			t.Fatalf("Failed to decode test string using i2pB64enc: '%v'", err)
		}

		if !bytes.Equal(decodedBytes, directlyDecoded) {
			t.Errorf("Mismatch between ToBytes result and direct decoding. ToBytes len: '%d', Direct decoding len: '%d'", len(decodedBytes), len(directlyDecoded))
		}
	})
}

/*
	func removeNewlines(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", ""), "\n", "")
	}
*/
func Test_KeyGenerationAndHandling(t *testing.T) {
	// Generate new keys
	keys, err := NewDestination()
	if err != nil {
		t.Fatalf("Failed to generate new I2P keys: %v", err)
	}
	t.Run("LoadKeysIncompat", func(t *testing.T) {
		//extract keys
		addr := keys.Address
		fmt.Println(addr)

		//both := removeNewlines(keys.Both)
		both := keys.Both
		fmt.Println(both)

		//FORMAT TO LOAD: (Address, Both)
		addrload := addr.Base64() + "\n" + both

		r := strings.NewReader(addrload)
		loadedKeys, err := LoadKeysIncompat(r)
		if err != nil {
			t.Fatalf("LoadKeysIncompat failed: %v", err)
		}

		if loadedKeys.Address != keys.Address {
			//fmt.Printf("loadedKeys.Address md5hash: '%s'\n keys.Address md5hash: '%s'\n", getMD5Hash(string(loadedKeys.Address)), getMD5Hash(string(keys.Address)))
			t.Errorf("LoadKeysIncompat returned incorrect address. Got '%s', want '%s'", loadedKeys.Address, keys.Address)

		}
		if loadedKeys.Both != keys.Both {
			t.Errorf("LoadKeysIncompat returned incorrect pair. Got '%s'\nwant '%s'\n", loadedKeys.Both, keys.Both)
			/*
				if loadedKeys.Both == removeNewlines(keys.Both) {
					fmt.Println("However, both pairs are correct if newline is removed in generated keys.")
				}
			*/
		}

	})

	expected := keys.Address.Base64() + "\n" + keys.Both

	t.Run("StoreKeysIncompat", func(t *testing.T) {
		var buf bytes.Buffer
		err := StoreKeysIncompat(*keys, &buf)
		if err != nil {
			t.Fatalf("StoreKeysIncompat failed: '%v'", err)
		}
		if buf.String() != expected {
			t.Errorf("StoreKeysIncompat wrote incorrect data. Got '%s', want '%s'", buf.String(), expected)
		}
	})

	t.Run("StoreKeys", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "test_keys_")
		if err != nil {
			t.Fatalf("Failed to create temp directory: '%v'", err)
		}
		defer os.RemoveAll(tmpDir)
		tmpFilePath := filepath.Join(tmpDir, "test_keys.txt")

		err = StoreKeys(*keys, tmpFilePath)
		if err != nil {
			t.Fatalf("StoreKeys failed: '%v'", err)
		}

		content, err := ioutil.ReadFile(tmpFilePath)
		if err != nil {
			t.Fatalf("Failed to read temp file: '%v'", err)
		}

		if string(content) != expected {
			t.Errorf("StoreKeys wrote incorrect data. Got '%s', want '%s'", string(content), expected)
		}
	})
}
