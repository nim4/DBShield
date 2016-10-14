package dbms

import (
	"bytes"
	"crypto/tls"
	"net"

	"github.com/nim4/DBShield/dbshield/logger"
)

//Postgres DBMS
type Postgres struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   string
	username    string
}

//SetCertificate to use if client asks for SSL
func (p *Postgres) SetCertificate(crt, key string) (err error) {
	p.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
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
func (p *Postgres) Handler() error {
	defer handlePanic()

	success, err := p.handleLogin()
	if err != nil {
		return err
	}
	if !success {
		logger.Warning("Login failed")
		return nil
	}
	return nil
}

func (p *Postgres) handleLogin() (success bool, err error) {
	//Receive Greeting
	buf, err := readPacket(p.client)
	if err != nil {
		return
	}

	//Send Greeting
	_, err = p.server.Write(buf)
	if err != nil {
		return
	}

	//Receive Greeting
	buf, err = readPacket(p.server)
	if err != nil {
		return
	}

	//Send Greeting
	_, err = p.client.Write(buf)
	if err != nil {
		return
	}

	//Receive username and database name
	buf, err = readPacket(p.client)
	if err != nil {
		return
	}
	data := buf[13:]
	nullByteIndex := bytes.IndexByte(data[13:], 0x00)
	logger.Infof("Username: %s", data[8:nullByteIndex+1])
	_, err = p.server.Write(buf)
	if err != nil {
		return
	}
	return
}
