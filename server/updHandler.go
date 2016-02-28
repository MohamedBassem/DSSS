package server

import (
	"fmt"
	"net"
)

func initUDP(agentUDPPort int) {

	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", agentUDPPort))
	if err != nil {
		panic(err)
	}

	ln, err := net.ListenUDP("udp", udpAddress)
	logger.Printf("Server is serving UDP pings on port 0.0.0.0:%v\n", agentUDPPort)
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 4)
		_, from, err := ln.ReadFromUDP(buf)
		if err == nil && string(buf) == "PING" {
			ln.WriteToUDP([]byte("PONG"), from)
		}
	}

}
