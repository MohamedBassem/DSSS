package agent

import (
	"log"
	"net"
	"strings"

	"github.com/MohamedBassem/DSSS/tcpreadwriter"
)

var Logger *log.Logger
var readWrite *tcpreadwriter.TCPReadWriter
var server string

func Main(l *log.Logger, addr string) {

	Logger = l
	server = addr
	if server == "" {
		Logger.Fatalln("You should specify the server's address (e.g. localhost:8082)")
	}
	InitTCPCon(server)

}

func startAgent(con *net.TCPConn) {

	readWrite = tcpreadwriter.New(con)

	for {
		msg, err := readWrite.ReadMessage()
		if err != nil {
			panic(err)
		}
		Logger.Println(msg)

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
	err := readWrite.WriteMessage("PONG")
	if err != nil {
		panic(err)
	}
}

func whoHas(arr []string) {
	exists := HasHash(arr[1])
	msg := "0"
	if exists {
		Logger.Println("Hash found")
		msg = "1"
	}
	err := readWrite.WriteMessage(msg)
	if err != nil {
		panic(err)
	}
}

func upload(arr []string) {
	Store(arr[1], strings.Join(arr[2:], " "))
	err := readWrite.WriteMessage("OK")
	if err != nil {
		panic(err)
	}
}

func download(arr []string) {
	cnt := Fetch(arr[1])
	err := readWrite.WriteMessage(cnt)
	if err != nil {
		panic(err)
	}
	Logger.Println("Content Sent")
}

func InitTCPCon(servAddr string) {

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
