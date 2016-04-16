package agent

import (
	"github.com/MohamedBassem/DSSS/tcpreadwriter"
	"net"
	"strings"
	"log"
)

const (
  server = "localhost:8082"
)


var Logger *log.Logger
var readWrite *tcpreadwriter.TCPReadWriter


func Main(l *log.Logger) {

	Logger = l
	InitTCPCon(server)	

}


func startAgent(con *net.TCPConn) {

	readWrite = tcpreadwriter.New(con)
	msg, err:= readWrite.ReadMessage()
	
	if err != nil {
		panic(err)
	}

	id := strings.Split(msg, " ")[1]
	Logger.Println(id)

	msg, err = readWrite.ReadMessage()

	for true {
		msg, err = readWrite.ReadMessage()
		if err != nil {
			panic(err)
		}
	
		arr := strings.Split(msg, " ")
		cmd := arr[0]

		if cmd == "PING" {
			ping(arr)		
		} else if cmd == "WHO_HAS" {
			whoHas(arr)
		} else if cmd == "UPLOAD" {
			upload(arr)
		} else if cmd == "DOWNLOAD" {
			download(arr)
		}
	}

}

func ping(arr []string) {
	Logger.Println("PING")
	err := readWrite.WriteMessage("PONG")
	if err != nil {
		panic(err)
	}
}

func whoHas(arr []string) {
	exists := HasHash(arr[1])
	msg := "0"
	if exists {
		msg = "1"	
	}
	err := readWrite.WriteMessage(msg)
	if err != nil {
		panic(err)
	}
}


func upload(arr []string) {
	Store(arr[1], strings.Join(arr[2:], ","))
	err := readWrite.WriteMessage("OK")
	if err != nil {
		panic(err)
	}
}

func download(arr []string) {
	cnt := Fetch(arr[1])
	err := readWrite.WriteMessage(cnt)
	if err != nil { panic(err) }
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

