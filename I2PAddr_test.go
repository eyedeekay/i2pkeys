package i2pkeys

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const (
	yoursam               = "127.0.0.1:7656"
	validShortenedI2PAddr = "i2p-projekt.i2p"
	validI2PAddrB32       = "udhdrtrcetjm5sxzskjyr5ztpeszydbh4dpl3pl4utgqqw2v4jna.b32.i2p"
	validI2PAddrB64       = "8ZAW~KzGFMUEj0pdchy6GQOOZbuzbqpWtiApEj8LHy2~O~58XKxRrA43cA23a9oDpNZDqWhRWEtehSnX5NoCwJcXWWdO1ksKEUim6cQLP-VpQyuZTIIqwSADwgoe6ikxZG0NGvy5FijgxF4EW9zg39nhUNKRejYNHhOBZKIX38qYyXoB8XCVJybKg89aMMPsCT884F0CLBKbHeYhpYGmhE4YW~aV21c5pebivvxeJPWuTBAOmYxAIgJE3fFU-fucQn9YyGUFa8F3t-0Vco-9qVNSEWfgrdXOdKT6orr3sfssiKo3ybRWdTpxycZ6wB4qHWgTSU5A-gOA3ACTCMZBsASN3W5cz6GRZCspQ0HNu~R~nJ8V06Mmw~iVYOu5lDvipmG6-dJky6XRxCedczxMM1GWFoieQ8Ysfuxq-j8keEtaYmyUQme6TcviCEvQsxyVirr~dTC-F8aZ~y2AlG5IJz5KD02nO6TRkI2fgjHhv9OZ9nskh-I2jxAzFP6Is1kyAAAA"
)

func Test_Basic(t *testing.T) {
	fmt.Println("Test_Basic")
	fmt.Println("\tAttaching to SAM at " + yoursam)
	keys, err := NewDestination()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(keys.String())
}

func Test_Basic_Lookup(t *testing.T) {
	fmt.Println("Test_Basic")
	fmt.Println("\tAttaching to SAM at " + yoursam)
	keys, err := Lookup("idk.i2p")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	fmt.Println(keys.String())
}

func Test_NewI2PAddrFromString(t *testing.T) {
	addr, err := NewI2PAddrFromString(validI2PAddrB64)

	if err != nil {
		t.Errorf("NewI2PAddrFromString failed for valid address: %v", err)
		t.Fail()
		return
	}

	if addr.Base64() != validI2PAddrB64 {
		t.Errorf("NewI2PAddrFromString returned incorrect address. Got %s, want %s", addr.Base64(), validI2PAddrB64)
		t.Fail()
		return
	}

	invalidAddr := "not-a-valid-address"
	_, err = NewI2PAddrFromString(invalidAddr)

	if err == nil {
		t.Error("NewI2PAddrFromString should have failed for invalid address")
		t.Fail()
		return
	}
	return
}

func Test_I2PAddrBase32(t *testing.T) {
	addr := I2PAddr(validI2PAddrB64)
	base32 := addr.Base32()

	if !strings.HasSuffix(base32, ".b32.i2p") {
		t.Errorf("Base32 address should end with .b32.i2p, got %s", base32)
		t.Fail()
		return
	}

	if len(base32) != 60 {
		t.Errorf("Base32 address should be 60 characters long, got %d", len(base32))
		t.Fail()
		return
	}
}

func Test_LoadKeysIncompat(t *testing.T) {
	testKeys := "testpub\ntestpriv"
	r := strings.NewReader(testKeys)

	keys, err := LoadKeysIncompat(r)

	if err != nil {
		t.Errorf("LoadKeysIncompat failed: %v", err)
		t.Fail()
		return
	}

	if keys.Address != I2PAddr("testpub") {
		t.Errorf("LoadKeysIncompat returned incorrect public key. Got %s, want %s", keys.Address, "testpub")
		t.Fail()
		return
	}

	if keys.Both != "testpriv" {
		t.Errorf("LoadKeysIncompat returned incorrect private key. Got %s, want %s", keys.Both, "testpriv")
		t.Fail()
		return
	}
}

func Test_StoreKeysIncompat(t *testing.T) {
	keys := I2PKeys{Address: I2PAddr("testpub"), Both: "testpriv"}
	var buf bytes.Buffer

	err := StoreKeysIncompat(keys, &buf)

	if err != nil {
		t.Errorf("StoreKeysIncompat failed: %v", err)
		t.Fail()
		return
	}

	expected := "testpub\ntestpriv"
	if buf.String() != expected {
		t.Errorf("StoreKeysIncompat wrote incorrect data. Got %s, want %s", buf.String(), expected)
		t.Fail()
		return
	}
}

func Test_DestHashFromString(t *testing.T) {

	hash, err := DestHashFromString(validI2PAddrB32)

	if err != nil {
		t.Errorf("DestHashFromString failed for valid hash: %v", err)
		t.Fail()
		return
	}

	if hash.String() != validI2PAddrB32 {
		t.Errorf("DestHashFromString returned incorrect hash. Got %s, want %s", hash.String(), validI2PAddrB32)
		t.Fail()
		return
	}

	invalidHash := "not-a-valid-hash"
	_, err = DestHashFromString(invalidHash)

	if err == nil {
		t.Error("DestHashFromString should have failed for invalid hash")
		t.Fail()
		return
	}
}

func Test_I2PAddrToBytes(t *testing.T) {
	addr := I2PAddr(validI2PAddrB64)

	decodedBytes, err := addr.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	//encode back to i2p-base64
	encodedString := i2pB64enc.EncodeToString(decodedBytes)

	if encodedString != validI2PAddrB64 {
		t.Errorf("Round-trip encoding/decoding failed. Got %s, want %s", encodedString, validI2PAddrB64)
	}

	//try decoding directly using i2p-base64
	directlyDecoded, err := i2pB64enc.DecodeString(validI2PAddrB64)
	if err != nil {
		t.Fatalf("Failed to decode test string using i2pB64enc: %v", err)
	}

	if !bytes.Equal(decodedBytes, directlyDecoded) {
		t.Errorf("Mismatch between ToBytes result and direct decoding. ToBytes len: %d, Direct decoding len: %d", len(decodedBytes), len(directlyDecoded))
	}
}
