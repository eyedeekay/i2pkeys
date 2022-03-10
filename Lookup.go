package i2pkeys

import "net"

type SimpleLookup struct {
}

func Lookup(i2paddr string) net.Addr {
	return nil
}

func LookupI2P(i2paddr string) *I2PAddr {
	return nil
}
