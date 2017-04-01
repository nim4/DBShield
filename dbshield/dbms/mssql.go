package dbms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

const maxMSSQLPayloadLen = 4096

const (
	mssqlEncryptionAvailableAndOff = iota
	mssqlEncryptionAvailableAndOn
	mssqlEncryptionNotAvailable
	mssqlEncryptionRequired
)

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
	return 1433
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
		if err != nil || len(buf) < 8 {
			return err
		}

		switch buf[0] {
		case 0x01: //SQL batch
			query := buf[8:]
			context := sql.QueryContext{
				Query:    query,
				Database: m.currentDB,
				User:     m.username,
				Client:   remoteAddrToIP(m.client.RemoteAddr()),
				Time:     time.Now(),
			}
			processContext(context)
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

	//Receive PreLogin Request
	buf, err := m.reader(m.client)
	if err != nil {
		return
	}

	if buf[0] != 0x12 {
		err = errors.New("packet is not PRELOGIN")
		return
	}

	//Send PreLogin to server
	_, err = m.server.Write(buf)
	if err != nil {
		return
	}

	//Receive PRELOGIN response
	buf, err = m.reader(m.server)
	if err != nil {
		return
	}

	if buf[0] != 0x4 {
		err = errors.New("packet is not PRELOGIN response")
		return
	}

	//Send PreLogin to server
	_, err = m.client.Write(buf)
	if err != nil {
		return
	}

	//Set data to beginning of the prelogin message
	data := buf[8:]

	var encryption byte

	//Lookup Encryption
	for i := 0; i < len(data); i += 5 {
		switch data[i] {
		case 0x01: //Encryption
			encryption = data[int(data[i+1])]
			break
		case 0xff: //Terminator
			break
		}
	}
	logger.Debugf("Encryption: %v", encryption)

	// buf, err = m.reader(m.client)
	// if err != nil {
	// 	return
	// }
	//
	//m.username, m.currentDB = MSSQLGetUsernameDB(buf)

	for {
		//Receive PreLogin Request
		err = readWrite(m.client, m.server, m.reader)
		if err != nil {
			return
		}

		//Receive PRELOGIN response
		buf, err = m.reader(m.server)
		if err != nil {
			return
		}

		//Send PreLogin to server
		_, err = m.client.Write(buf)
		if err != nil {
			return
		}

		if buf[0] == 0x4 {
			break
		}
	}
	success = true
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
	//TODO: Extract Username and db name
	fmt.Printf("%x", data)
	return
}
