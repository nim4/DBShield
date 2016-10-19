package dbms

import (
	"bytes"
	"crypto/tls"
	"net"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

//MySQL DBMS
type MySQL struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   string
	username    string
	reader      func(net.Conn) ([]byte, error)
}

//SetCertificate to use if client asks for SSL
func (m *MySQL) SetCertificate(crt, key string) (err error) {
	m.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

//SetReader function for sockets IO
func (m *MySQL) SetReader(f func(net.Conn) ([]byte, error)) {
	m.reader = f
}

//SetSockets for dbms (client and server sockets)
func (m *MySQL) SetSockets(c, s net.Conn) {
	m.client = c
	m.server = s
}

//Close sockets
func (m *MySQL) Close() {
	m.client.Close()
	m.server.Close()
}

//DefaultPort of the DBMS
func (m *MySQL) DefaultPort() uint {
	return 3306
}

//Handler gets incoming requests
func (m *MySQL) Handler() error {
	defer handlePanic()

	success, err := m.handleLogin()
	if err != nil {
		return err
	}
	if !success {
		logger.Warning("Login failed")
		return nil
	}
	for {
		var buf []byte
		buf, err = m.reader(m.client)
		if err != nil {
			return err
		}
		data := buf[4:]

		switch data[0] {
		case 0x01: //Quit
			return nil
		case 0x02: //UseDB
			m.currentDB = string(data[1:])
			logger.Infof("Using database: %v", m.currentDB)
		case 0x03: //Query
			query := data[1:]
			logger.Infof("Query: %s", query)
			context := sql.QueryContext{
				Query:    string(query),
				Database: m.currentDB,
				User:     m.username,
				Client:   remoteAddrToIP(m.client.RemoteAddr()),
				Time:     time.Now().Unix(),
			}
			processContext(context)
			// case 0x04: //Show fields
			// 	logger.Debugf("Show fields: %s", data[1:])
			// default:
			// 	logger.Debugf("Unknown Data[0]: %x", data[0])
		}

		//Send query/request to server
		_, err = m.server.Write(buf)
		if err != nil {
			return err
		}
		//Recive response
		buf, err = m.reader(m.server)
		if err != nil {
			return err
		}

		//Send response to client
		_, err = m.client.Write(buf)
		if err != nil {
			return err
		}
	}
}

func (m *MySQL) handleLogin() (success bool, err error) {
	//Receive Server Greeting
	buf, err := m.reader(m.server)
	if err != nil {
		return
	}
	/* Extra Info
	data := buf[4:]
	nullByteIndex := bytes.IndexByte(data[1:], 0x00)
	logger.Infof("Version: %s", data[1:nullByteIndex+1])
	pos := 1 + nullByteIndex + 1 + 4
	// first part of the password cipher [8 bytes]
	cipher := make([]byte, 20)
	copy(cipher, data[pos:pos+8])
	pos += 8 + 1 + 2 + 1 + 2 + 2 + 1 + 10
	cipher = append(cipher[:8], data[pos:pos+12]...)
	logger.Debugf("Cipher: 0x%x", cipher)
	if err != nil {
		return
	}*/

	//Send Server Greeting to client
	_, err = m.client.Write(buf)
	if err != nil {
		return
	}

	//Receive Login Request
	buf, err = m.reader(m.client)
	if err != nil {
		return
	}
	data := buf[4:]

	m.username, _, m.currentDB = getUsernameHashDB(data)

	//check if ssl is required
	ssl := (data[1] & 0x08) == 0x08

	//Send Login Request
	_, err = m.server.Write(buf)
	if err != nil {
		return
	}
	if ssl {
		m.client, m.server, err = turnSSL(m.client, m.server, m.certificate)

		buf, err = m.reader(m.client)
		if err != nil {
			return
		}
		data = buf[4:]
		m.username, _, m.currentDB = getUsernameHashDB(data)

		//Send Login Request
		_, err = m.server.Write(buf)
		if err != nil {
			return
		}
	}
	logger.Debugf("SSL bit: %v", ssl)

	if len(m.currentDB) != 0 { //db Selected
		//Receive OK
		buf, err = m.reader(m.server)
		if err != nil {
			return
		}
	} else {
		//Receive Auth Switch Request
		buf, err = m.reader(m.server)
		if err != nil {
			return
		}
		//Send Auth Switch Response
		_, err = m.client.Write(buf)
		if err != nil {
			return
		}
		//Receive Auth Switch Response
		buf, err = m.reader(m.client)
		if err != nil {
			return
		}

		//Send Auth Switch Response
		_, err = m.server.Write(buf)
		if err != nil {
			return
		}
		//Receive Response Status
		buf, err = m.reader(m.server)
		if err != nil {
			return
		}
	}

	if buf[5] != 0x15 {
		success = true
	}

	//Send Response Status
	_, err = m.client.Write(buf)
	if err != nil {
		return
	}
	return
}

func getUsernameHashDB(data []byte) (username string, hash []byte, db string) {
	if len(data) < 33 {
		return
	}
	pos := 32

	nullByteIndex := bytes.IndexByte(data[pos:], 0x00)
	username = string(data[pos : nullByteIndex+pos])
	logger.Infof("Username: %s", username)
	pos += nullByteIndex + 2
	hash = data[pos : pos+20]
	logger.Infof("Hash: %x", data[pos:pos+20])
	pos += 20
	nullByteIndex = bytes.IndexByte(data[pos:], 0x00)
	if nullByteIndex != 0 {
		db = string(data[pos : nullByteIndex+pos])
		logger.Infof("Database: %s", db)
	}
	return
}
