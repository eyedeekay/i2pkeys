package i2pkeys

import (
	"fmt"
	"net"
	"strings"
)

func Lookup(addr string) (*I2PAddr, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("HELLO VERSION MIN=3.1 MAX=3.1\n"))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < 1 {
		return nil, fmt.Errorf("no data received")
	}
	if strings.Contains(string(buf[:n]), "RESULT=OK") {
		_, err = conn.Write([]byte(fmt.Sprintf("NAMING LOOKUP NAME=%s\n", addr)))
		if err != nil {
			return nil, err
		}
		n, err = conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if n < 1 {
			return nil, fmt.Errorf("no destination data received")
		}
		value := strings.Split(string(buf[:n]), "VALUE=")[1]
		addr, err := NewI2PAddrFromString(value)
		if err != nil {
			return nil, err
		}
		return &addr, err
	}
	return nil, fmt.Errorf("no result received")
}
