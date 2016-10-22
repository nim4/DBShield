package dbms

import (
	"bytes"
	"net"
)

const (
	chunkSize = 256
)

//ReadPacket all available data from socket
func ReadPacket(conn net.Conn) ([]byte, error) {
	buf := &bytes.Buffer{}
	for {
		data := make([]byte, chunkSize)
		n, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		buf.Write(data[:n])
		if n != chunkSize {
			break
		}
	}
	return buf.Bytes(), nil
}
