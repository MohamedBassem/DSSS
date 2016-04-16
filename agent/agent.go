package agent

import (
	"github.com/MohamedBassem/DSSS/tcpreadwriter"
	"net"
	"log"
)

const (
	chunkSize                    = 100 // In bytes
  server											 = "localhost:8082"
)


var logger *log.Logger



func Main(l *log.Logger) {

	logger = l
	InitTCPCon(server)	

}

func startAgent(con *net.TCPConn) {

	readWrite := tcpreadwriter.New(con)
	msg, err:= readWrite.ReadMessage()
	
	if err != nil {
		panic(err)
	}
	
	logger.Println(msg)
}

func InitTCPCon(servAddr string){

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        panic(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        panic(err)
    }

	  startAgent(conn)

    conn.Close()


}

