package dbms

import (
	"bytes"
	"io"
	"net"
)

const (
	chunkSize = 4096
)

//ReadPacket all available data from socket
func ReadPacket(conn io.Reader) ([]byte, error) {
	data := make([]byte, chunkSize)
	buf := bytes.Buffer{}
	for {
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

func readWrite(src, dst net.Conn, reader func(io.Reader) ([]byte, error)) error {
	//Read result from server
	buf, err := reader(src)
	if err != nil {
		return err
	}

	//Send result to client
	_, err = dst.Write(buf)
	return err
}
