package i2pkeys

import (
	"fmt"
	"testing"
	//	"time"
)

const yoursam = "127.0.0.1:7656"

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
