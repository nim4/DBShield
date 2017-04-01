package dbms

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

const maxMSSQLPayloadLen = 1<<24 - 1

//MSSQL DBMS
type MSSQL struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

//SetCertificate to use if client asks for SSL
func (m *MSSQL) SetCertificate(crt, key string) (err error) {
	m.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

//SetReader function for sockets IO
func (m *MSSQL) SetReader(f func(io.Reader) ([]byte, error)) {
	m.reader = f
}

//SetSockets for dbms (client and server sockets)
func (m *MSSQL) SetSockets(c, s net.Conn) {
	m.client = c
	m.server = s
}

//Close sockets
func (m *MSSQL) Close() {
	defer handlePanic()
	m.client.Close()
	m.server.Close()
}

//DefaultPort of the DBMS
func (m *MSSQL) DefaultPort() uint {
	return 3306
}

//Handler gets incoming requests
func (m *MSSQL) Handler() error {
	defer handlePanic()
	defer m.Close()
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
		if err != nil || len(buf) < 5 {
			return err
		}
		data := buf[4:]

		switch data[0] {
		case 0x01: //Quit
			return nil
		case 0x02: //UseDB
			m.currentDB = data[1:]
			logger.Debugf("Using database: %v", m.currentDB)
		case 0x03: //Query
			query := data[1:]
			context := sql.QueryContext{
				Query:    query,
				Database: m.currentDB,
				User:     m.username,
				Client:   remoteAddrToIP(m.client.RemoteAddr()),
				Time:     time.Now(),
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
		err = readWrite(m.server, m.client, m.reader)
		if err != nil {
			return err
		}
	}
}

func (m *MSSQL) handleLogin() (success bool, err error) {
	//Receive Server Greeting
	err = readWrite(m.server, m.client, m.reader)
	if err != nil {
		return
	}

	//Receive Login Request
	buf, err := m.reader(m.client)
	if err != nil {
		return
	}
	data := buf[4:]

	m.username, m.currentDB = MSSQLGetUsernameDB(data)

	//check if ssl is required
	ssl := (data[1] & 0x08) == 0x08

	//Send Login Request
	_, err = m.server.Write(buf)
	if err != nil {
		return
	}
	if ssl {
		m.client, m.server, err = turnSSL(m.client, m.server, m.certificate)
		if err != nil {
			return
		}
		buf, err = m.reader(m.client)
		if err != nil {
			return
		}
		data = buf[4:]
		m.username, m.currentDB = MSSQLGetUsernameDB(data)

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
		err = readWrite(m.server, m.client, m.reader)
		if err != nil {
			return
		}
		//Receive Auth Switch Response
		err = readWrite(m.client, m.server, m.reader)
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

//buffer pool for MSSQLReadPacket
var mssqlDataPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, maxMSSQLPayloadLen)
	},
}

//MSSQLReadPacket handles reading mysql packets
func MSSQLReadPacket(src io.Reader) ([]byte, error) {
	data := mssqlDataPool.Get().([]byte)
	defer mssqlDataPool.Put(data)
	n, err := src.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if n < maxMSSQLPayloadLen || err == io.EOF {
		return data[:n], nil
	}

	buf, err := MSSQLReadPacket(src)
	if err != nil {
		return nil, err
	}
	return append(data, buf...), nil
}

//MSSQLGetUsernameDB parse packet and gets username and db name
func MSSQLGetUsernameDB(data []byte) (username, db []byte) {
	if len(data) < 33 {
		return
	}
	pos := 32

	nullByteIndex := bytes.IndexByte(data[pos:], 0x00)
	username = data[pos : nullByteIndex+pos]
	logger.Debugf("Username: %s", username)
	pos += nullByteIndex + 22
	nullByteIndex = bytes.IndexByte(data[pos:], 0x00)

	//Check if DB name is selected
	dbSelectedCheck := len(data) > nullByteIndex+pos+1

	if nullByteIndex != 0 && dbSelectedCheck {
		db = data[pos : nullByteIndex+pos]
		logger.Debugf("Database: %s", db)
	}
	return
}
