package dbms

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

//Postgres DBMS
type Postgres struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

//SetCertificate to use if client asks for SSL
func (p *Postgres) SetCertificate(crt, key string) (err error) {
	p.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

//SetReader function for sockets IO
func (p *Postgres) SetReader(f func(io.Reader) ([]byte, error)) {
	p.reader = f
}

//SetSockets for dbms (client and server sockets)
func (p *Postgres) SetSockets(c, s net.Conn) {
	defer handlePanic()
	p.client = c
	p.server = s
}

//Close sockets
func (p *Postgres) Close() {
	defer handlePanic()
	p.client.Close()
	p.server.Close()
}

//DefaultPort of the DBMS
func (p *Postgres) DefaultPort() uint {
	return 5432
}

//Handler gets incoming requests
func (p *Postgres) Handler() (err error) {
	//defer handlePanic()
	defer p.Close()

	success, err := p.handleLogin()
	if err != nil {
		return
	}
	if !success {
		logger.Warning("Login failed")
		return
	}
	for {
		var buf []byte
		//Read client request
		buf, err = p.reader(p.client)
		if err != nil {
			return
		}
		switch buf[0] {
		case 0x51: //Simple query
			context := sql.QueryContext{
				Query:    buf[5:],
				Database: p.currentDB,
				User:     p.username,
				Client:   remoteAddrToIP(p.client.RemoteAddr()),
				Time:     time.Now(),
			}
			processContext(context)

		case 0x58: //Terminate
			_, err = p.server.Write(buf)
			return
		}

		//Send request to server
		_, err = p.server.Write(buf)
		if err != nil {
			return
		}

		//Read server response
		buf, err = p.reader(p.server)
		if err != nil {
			return
		}

		//Send response to client
		_, err = p.client.Write(buf)
		if err != nil {
			return
		}

		switch buf[0] {
		case 0x45: //Error
			buf, err = p.reader(p.server)
			if err != nil {
				return
			}
			_, err = p.client.Write(buf)
			if err != nil {
				return
			}
		}
	}
}

func (p *Postgres) handleLogin() (success bool, err error) {
	//Receive Greeting
	err = readWrite(p.client, p.server, p.reader)
	if err != nil {
		return
	}

	//Receive Greeting
	buf, err := p.reader(p.server)
	if err != nil {
		return
	}
	ssl := buf[0] == 0x53

	//Send Greeting
	_, err = p.client.Write(buf)
	if err != nil {
		return
	}

	if ssl {
		p.client, p.server, err = turnSSL(p.client, p.server, p.certificate)
		if err != nil {
			return
		}
	}

	//Receive username and database name
	buf, err = p.reader(p.client)
	if err != nil {
		return
	}

	data := buf[8:]

	payload := make(map[string][]byte)
	for {
		//reading key
		nullByteIndex := bytes.IndexByte(data, 0x00)
		if nullByteIndex <= 0 {
			break
		}
		key := string(data[:nullByteIndex+1])

		//reading value
		data = data[nullByteIndex+1:]
		nullByteIndex = bytes.IndexByte(data, 0x00)
		if nullByteIndex <= 0 {
			break
		}
		payload[key] = data[:nullByteIndex+1]
		data = data[nullByteIndex+1:]
	}
	for key := range payload {
		logger.Debugf("%s: %s", strings.Title(key), payload[key])
	}
	p.username = payload["user"]
	p.currentDB = payload["database"]

	//Send username & dbname to server
	_, err = p.server.Write(buf)
	if err != nil {
		return
	}

	//Read authentication request from server
	err = readWrite(p.server, p.client, p.reader)
	if err != nil {
		return
	}

	//Read client password message
	err = readWrite(p.client, p.server, p.reader)
	if err != nil {
		return
	}

	//Read authtentication result from server
	buf, err = p.reader(p.server)
	if err != nil {
		return
	}
	data = buf[5:9]
	if binary.BigEndian.Uint32(data) == 0 {
		success = true
	}
	//Send authtentication result to client
	_, err = p.client.Write(buf)
	return
}
