package server

import "log"

const (
	agentIdLength     = 30
	replicationFactor = 1
)

var logger *log.Logger

func Main(agentTCPPort, clientHTTPPort int, l *log.Logger) {

	logger = l

	go initHTTP(clientHTTPPort)
	go initTCP(agentTCPPort)

	<-make(chan struct{}, 0)
}
