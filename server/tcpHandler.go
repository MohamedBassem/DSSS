package server

import (
	"fmt"
	"net"

	"github.com/MohamedBassem/DSSS/tcpreadwriter"
)

type query struct {
	text     string
	response chan<- response
}

type response struct {
	text string
	err  error
}

type agent struct {
	id         string
	tcpAddress *net.TCPAddr
	udpAddress *net.UDPAddr

	queries chan query
}

func handleTCPConnection(_conn *net.TCPConn) {
	_conn.SetKeepAlive(true)

	readWriter := tcpreadwriter.New(_conn)

	// Each agent has an ID to relate its TCP connection with its UDP connection
	agent := agent{
		id:         generateRandomString(agentIdLength),
		tcpAddress: _conn.RemoteAddr().(*net.TCPAddr),
		queries:    make(chan query, 100),
	}

	connectedAgents.put(agent.id, &agent)
	defer func() {
		_conn.Close()
		logger.Printf("Client %v disconnected from the TCP endpoint\n", agent.tcpAddress)
		connectedAgents.delete(agent.id)
	}()

	logger.Printf("Client %v connected to the TCP endpoint with id %v\n", agent.tcpAddress, agent.id)

	err := readWriter.WriteMessage(fmt.Sprintf("HI %v", agent.id))
	if err != nil {
		return
	}

	for {
		select {
		case query := <-agent.queries:
			err := readWriter.WriteMessage(query.text)
			if err != nil {
				query.response <- response{text: "", err: err}
				return
			}
			// TODO : Add Read Timeout
			message, err := readWriter.ReadMessage()
			query.response <- response{text: message, err: err}
			if err != nil {
				return
			}
			//case <-time.After(time.Second * 3):
			//err := readWriter.WriteMessage("PING")
			//if err != nil {
			//return
			//}
			//// TODO : Add Read Timeout
			//message, err := readWriter.ReadMessage()
			//if err != nil || message != "PONG" {
			//return
			//}
		}
	}

}

func initTCP(agentTCPPort int) {

	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", agentTCPPort))
	if err != nil {
		panic(err)
	}

	server, err := net.ListenTCP("tcp", tcpAddress)
	logger.Printf("Server is serving TCP on port 0.0.0.0:%v\n", agentTCPPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}
		go handleTCPConnection(conn)
	}

}
