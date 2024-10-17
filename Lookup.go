package i2pkeys

import (
	"fmt"
	"net"
	"strings"
)

func Lookup(addr string) (*I2PAddr, error) {
	log.WithField("addr", addr).Debug("Starting Lookup")
	conn, err := net.Dial("tcp", "127.0.0.1:7656")
	if err != nil {
		log.Error("Failed to connect to SAM bridge")
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("HELLO VERSION MIN=3.1 MAX=3.1\n"))
	if err != nil {
		log.Error("Failed to write HELLO VERSION")
		return nil, err
	}
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("Failed to read HELLO VERSION response")
		return nil, err
	}
	if n < 1 {
		log.Error("no data received")
		return nil, fmt.Errorf("no data received")
	}

	response := string(buf[:n])
	log.WithField("response", response).Debug("Received HELLO response")

	if strings.Contains(string(buf[:n]), "RESULT=OK") {
		_, err = conn.Write([]byte(fmt.Sprintf("NAMING LOOKUP NAME=%s\n", addr)))
		if err != nil {
			log.Error("Failed to write NAMING LOOKUP command")
			return nil, err
		}
		n, err = conn.Read(buf)
		if err != nil {
			log.Error("Failed to read NAMING LOOKUP response")
			return nil, err
		}
		if n < 1 {
			return nil, fmt.Errorf("no destination data received")
		}
		value := strings.Split(string(buf[:n]), "VALUE=")[1]
		addr, err := NewI2PAddrFromString(value)
		if err != nil {
			log.Error("Failed to parse I2P address from lookup response")
			return nil, err
		}
		log.WithField("addr", addr).Debug("Successfully resolved I2P address")
		return &addr, err
	}
	log.Error("no RESULT=OK received in HELLO response")
	return nil, fmt.Errorf("no result received")
}
