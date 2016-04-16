package agent

import (
	"github.com/MohamedBassem/DSSS/tcpreadwriter"
	"net"
	"strings"
	"log"
)

const (
  server											 = "localhost:8082"
)


var logger *log.Logger
var readWrite *tcpreadwriter.TCPReadWriter


func Main(l *log.Logger) {

	logger = l
	InitTCPCon(server)	

}

func startAgent(con *net.TCPConn) {

	readWrite = tcpreadwriter.New(con)
	msg, err:= readWrite.ReadMessage()
	
	if err != nil {
		panic(err)
	}

	id := strings.Split(msg, " ")[1]
	logger.Println(id)

	msg, err = readWrite.ReadMessage()

	for true {
		msg, err = readWrite.ReadMessage()
		if err != nil {
			panic(err)
		}
	
		arr := strings.Split(msg, " ")
		cmd := arr[0]

		if cmd == "PING" {
			ping()		
		} else if cmd == "WHO_HAS" {
			whoHas()
		} else if cmd == "UPLOAD" {
			upload()
		}
	}

}

func ping() {
	logger.Println("PING")
	err := readWrite.WriteMessage("PONG")
	if err != nil {
		panic(err)
	}
}

func whoHas() {

}

func upload() {

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

