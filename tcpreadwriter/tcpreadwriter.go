package tcpreadwriter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type TCPReadWriter struct {
	readerWriter *bufio.ReadWriter
}

func New(conn *net.TCPConn) *TCPReadWriter {

	bufio := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	return &TCPReadWriter{
		readerWriter: bufio,
	}
}

func (t *TCPReadWriter) WriteMessage(message string) error {

	messageBytes := []byte(message)

	var n int32 = int32(len(messageBytes))
	bs := new(bytes.Buffer)
	binary.Write(bs, binary.LittleEndian, n)

	fullMessageBytes := append(bs.Bytes(), messageBytes...)
	_, err := t.readerWriter.Write(fullMessageBytes)
	return err
}

func (t *TCPReadWriter) ReadMessage() (string, error) {

	messageLengthBytes := make([]byte, 4)

	_, err := io.ReadFull(t.readerWriter, messageLengthBytes)
	if err != nil {
		return "", err
	}

	buf := bytes.NewReader(messageLengthBytes)

	var messageLength int32
	err = binary.Read(buf, binary.LittleEndian, &messageLength)
	if err != nil {
		return "", err
	}

	fullMessageBytes := make([]byte, messageLength)
	_, err = io.ReadFull(t.readerWriter, fullMessageBytes)

	if err != nil {
		return "", err
	}
	return string(fullMessageBytes), nil
}
