package dbms

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

//Postgres DBMS
type Postgres struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   string
	username    string
	reader      func(net.Conn) ([]byte, error)
}

//SetCertificate to use if client asks for SSL
func (p *Postgres) SetCertificate(crt, key string) (err error) {
	p.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

//SetReader function for sockets IO
func (p *Postgres) SetReader(f func(net.Conn) ([]byte, error)) {
	p.reader = f
}

//SetSockets for dbms (client and server sockets)
func (p *Postgres) SetSockets(c, s net.Conn) {
	p.client = c
	p.server = s
}

//Close sockets
func (p *Postgres) Close() {
	p.client.Close()
	p.server.Close()
}

//DefaultPort of the DBMS
func (p *Postgres) DefaultPort() uint {
	return 5432
}

//Handler gets incoming requests
func (p *Postgres) Handler() (err error) {
	defer handlePanic()
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
			query := buf[5:]
			fmt.Printf("Query: %s", query)
			logger.Infof("Query: %s", query)
			context := sql.QueryContext{
				Query:    string(query),
				Database: p.currentDB,
				User:     p.username,
				Client:   remoteAddrToIP(p.client.RemoteAddr()),
				Time:     time.Now().Unix(),
			}
			if config.Config.Learning {
				go training.AddToTrainingSet(context)
			} else {
				if config.Config.ActionFunc != nil && !training.CheckQuery(context) {
					return config.Config.ActionFunc(p.client)
				}
			}

		case 0x58: //Terminate
			_, err = p.server.Write(buf)
			return
		}

		//Send request to server
		_, err = p.server.Write(buf)
		if err != nil {
			return
		}

		//Read result from server
		buf, err = p.reader(p.server)
		if err != nil {
			return
		}

		//Send result to client
		_, err = p.client.Write(buf)
		if err != nil {
			return
		}
	}
}

func (p *Postgres) handleLogin() (success bool, err error) {
	//Receive Greeting
	buf, err := p.reader(p.client)
	if err != nil {
		return
	}

	//Send Greeting
	_, err = p.server.Write(buf)
	if err != nil {
		return
	}

	//Receive Greeting
	buf, err = p.reader(p.server)
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
		logger.Info("SSL connection")
		tlsConnClient := tls.Server(p.client, &tls.Config{
			Certificates:       []tls.Certificate{p.certificate},
			InsecureSkipVerify: true,
		})
		if err = tlsConnClient.Handshake(); err != nil {
			return
		}
		p.client = tlsConnClient
		logger.Debug("Client handshake done")

		tlsConnServer := tls.Client(p.server, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err = tlsConnServer.Handshake(); err != nil {
			return
		}
		p.server = tlsConnServer
		logger.Debug("Server handshake done")
	}

	//Receive username and database name
	buf, err = p.reader(p.client)
	if err != nil {
		return
	}

	data := buf[8:]

	payload := make(map[string]string)
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
		payload[key] = string(data[:nullByteIndex+1])
		data = data[nullByteIndex+1:]
	}
	for key := range payload {
		logger.Infof("%s: %s", strings.Title(key), payload[key])
	}
	p.username = payload["user"]
	p.currentDB = payload["database"]

	//Send username & dbname to server
	_, err = p.server.Write(buf)
	if err != nil {
		return
	}

	//Read authentication request from server
	buf, err = p.reader(p.server)
	if err != nil {
		return
	}

	//Send response to client
	_, err = p.client.Write(buf)
	if err != nil {
		return
	}

	//Read client password message
	buf, err = p.reader(p.client)
	if err != nil {
		return
	}

	//Send password to server
	_, err = p.server.Write(buf)
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
	if err != nil {
		return
	}
	return
}
