package server

import "log"

var logger *log.Logger

func Main(agentTCPPort, agentUDPPort, clientHTTPPort int, l *log.Logger) {

	logger = l

	go initHTTP(clientHTTPPort)

	<-make(chan struct{}, 0)
}
