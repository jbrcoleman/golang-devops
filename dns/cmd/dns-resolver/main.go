package main

import (
	"fmt"
	"net"
	"dns/pkg/dns"
)

func main() {
	fmt.Printf("Starting DNS Server...\n")
	packetConnection, err := net.ListenPacket("udp", ":53")
	if err != nil {
		panic(err)
	}
	defer packetConnection.Close()
    for {
		buf := make([]byte, 512)
        bytesRead, addr, err := packetConnection.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Read error from %s: %s", addr.String(), err)
			continue
		}
		go dns.HandlePacket(packetConnection, addr, buf[:bytesRead])
	}
}

