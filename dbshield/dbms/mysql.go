package dbms

import (
	"bytes"
	"crypto/tls"
	"net"
	"time"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

//MySQL DBMS
type MySQL struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   string
	username    string
}

//SetCertificate to use if client asks for SSL
func (m *MySQL) SetCertificate(crt, key string) (err error) {
	m.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
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
		buf, err = readPacket(m.client)
		if err != nil {
			return err
		}
		queryRequest := false
		data := buf[4:]

		switch data[0] {
		case 0x01: //Quit
			return nil
		case 0x02: //UseDB
			m.currentDB = string(data[1:])
			logger.Infof("Using database: %v", m.currentDB)
		case 0x03: //Query
			queryRequest = true
			query := data[1:]
			logger.Infof("Query: %s", query)
			context := sql.QueryContext{
				Query:  string(query),
				User:   m.username,
				Client: remoteAddrToIP(m.client.RemoteAddr()),
				Time:   time.Now(),
			}
			if config.Config.Learning {
				go training.AddToTrainingSet(context)
			} else {
				if config.Config.Action != nil && !training.CheckQuery(context) {
					return config.Config.Action(m.client)
				}
			}
		case 0x04: //Show fields
			logger.Debugf("Show fields: %s", data[1:])
		default:
			logger.Debugf("Unknown Data[0]: %x", data[0])
		}

		//Send query/request to server
		_, err = m.server.Write(buf)
		if err != nil {
			return err
		}
		//Recive response
		buf, err = readPacket(m.server)
		if err != nil {
			return err
		}

		if queryRequest && buf[0] == 0x1 {
			orginalColumns := []string{}
			var db, orginalTable string
			data = buf[5:]
			//Get Columns
			for i := 0; uint(i) < uint(buf[4]); i++ {
				var pos uint = 4
				//Catalog
				_, n := pascalString(data[pos:])
				pos += n + 1

				//DB name
				db, n = pascalString(data[pos:])
				pos += n + 1

				//in query table name
				_, n = pascalString(data[pos:])
				pos += n + 1

				//orginal table name
				orginalTable, n = pascalString(data[pos:])
				pos += n + 1

				//in query name ("count(*)", ...)
				qName, n := pascalString(data[pos:])
				pos += n + 1

				//orginal name
				orginalName, n := pascalString(data[pos:])
				pos += n + 1

				//If name is empty use qName so user can still write Masking rule
				if len(orginalName) == 0 {
					orginalName = qName
				}

				orginalColumns = append(orginalColumns, orginalName)
				// Filler [uint8]
				// Charset [charset, collation uint8]
				// Length [uint32]
				pos += 1 + 2 + 4

				// Field type [uint8]
				//fieldType := data[pos]
				pos++
				// Flags [uint16]
				//flags := uint16(binary.LittleEndian.Uint16(data[pos : pos+2]))
				pos += 2

				// Decimals [uint8]
				//decimals := data[pos]

				data = data[threeByteBigEndianToInt(data[:3])+4:]

			}
			logger.Debugf("Columns: %v", orginalColumns)
			//Get rows
			if len(orginalColumns) > 0 {
				if (data[0] == 0x05) && (len(data) > 4 || data[4] == 0xfe) { //Skip first  segment
					data = data[9:] //0x05 + 4
				}
			RowLoop:
				for {
					var rowLen = threeByteBigEndianToInt(data[:3]) + 4
					var pos uint = 4
					for _, col := range orginalColumns {
						if rowLen == pos {
							break
						}
						if (data[pos] == 0xff) || (data[pos] == 0xfe) { //}&& (data[pos+1] == 0) && (data[pos+2] == 0)) {
							break RowLoop
						} else if data[pos] == 0xfb {
							pos++
							continue
						}
						val, n := pascalString(data[pos:])
						logger.Debugf("Val: %s", val)
						pos++
						//Key to find masking rule
						key := db + "." + orginalTable + "." + col
						logger.Debugf("Key: %s", key)

						mask, ok := getMask(key, []byte(val))
						if ok {
							copy(data[pos:], mask)
						}

						pos += n

					}
					data = data[threeByteBigEndianToInt(data[:3])+4:]
				}
			}
		}

		_, err = m.client.Write(buf)
		if err != nil {
			return err
		}
	}
}

func (m *MySQL) handleLogin() (success bool, err error) {
	//Receive Server Greeting
	buf, err := readPacket(m.server)
	if err != nil {
		return
	}
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
	}

	//Send Server Greeting to client
	_, err = m.client.Write(buf)
	if err != nil {
		return
	}

	//Receive Login Request
	buf, err = readPacket(m.client)
	if err != nil {
		return
	}
	data = buf[4:]

	m.username, _, m.currentDB = getUsernameHashDB(data)

	//check if ssl is required
	ssl := (data[1] & 0x08) == 0x08

	//Send Login Request
	_, err = m.server.Write(buf)
	if err != nil {
		return
	}
	if ssl {
		logger.Info("SSL connection")
		tlsConnClient := tls.Server(m.client, &tls.Config{
			Certificates:       []tls.Certificate{m.certificate},
			InsecureSkipVerify: true,
		})
		if err = tlsConnClient.Handshake(); err != nil {
			return
		}
		m.client = tlsConnClient
		logger.Debug("Client handshake done")

		//Read TLS Hello
		buf, err = readPacket(m.client)
		if err != nil {
			return
		}
		data = buf[4:]
		m.username, _, m.currentDB = getUsernameHashDB(data)
		tlsConnServer := tls.Client(m.server, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err = tlsConnServer.Handshake(); err != nil {
			return
		}
		m.server = tlsConnServer
		logger.Debug("Server handshake done")

		//Send Login Request
		_, err = m.server.Write(buf)
		if err != nil {
			return
		}
	}
	logger.Debugf("SSL bit: %v", ssl)

	if len(m.currentDB) != 0 { //db Selected
		//Receive OK
		buf, err = readPacket(m.server)
		if err != nil {
			return
		}
	} else {
		//Receive Auth Switch Request
		buf, err = readPacket(m.server)
		if err != nil {
			return
		}
		//Send Auth Switch Response
		_, err = m.client.Write(buf)
		if err != nil {
			return
		}
		//Receive Auth Switch Response
		buf, err = readPacket(m.client)
		if err != nil {
			return
		}

		//Send Auth Switch Response
		_, err = m.server.Write(buf)
		if err != nil {
			return
		}
		//Receive Response Status
		buf, err = readPacket(m.server)
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
