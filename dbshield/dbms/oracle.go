package dbms

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"time"

	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

//Oracle DBMS
type Oracle struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

//SetCertificate to use if client asks for SSL
func (o *Oracle) SetCertificate(crt, key string) (err error) {
	o.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

//SetReader function for sockets IO
func (o *Oracle) SetReader(f func(io.Reader) ([]byte, error)) {
	o.reader = f
}

//SetSockets for dbms (client and server sockets)
func (o *Oracle) SetSockets(c, s net.Conn) {
	defer handlePanic()
	o.client = c
	o.server = s
}

//Close sockets
func (o *Oracle) Close() {
	defer handlePanic()
	o.client.Close()
	o.server.Close()
}

//DefaultPort of the DBMS
func (o *Oracle) DefaultPort() uint {
	return 1521
}

//Handler gets incoming requests
func (o *Oracle) Handler() error {
	defer handlePanic()
	defer o.Close()

	for {
		buf, err := o.readPacket(o.client)
		if err != nil {
			return err
		}
		var eof bool

		switch buf[4] { //Packet Type
		case 0x01: //Connect
			connectDataLen := int(buf[24])*256 + int(buf[25])
			connectData := buf[len(buf)-connectDataLen:]

			//Extracting Service name
			// FIXME: avoid string
			tmp1 := bytes.Split(connectData, []byte("SERVICE_NAME="))
			tmp2 := bytes.Split(tmp1[1], []byte{0x29}) // )
			o.currentDB = tmp2[0]

			logger.Debugf("Connect Data: %s", connectData)
			logger.Debugf("Service Name: %s", o.currentDB)
		case 0x06: //Data
			data := buf[8:]
			eof = data[1] == 0x40
			if !eof {
				payload := data[2:]
				if len(payload) > 16 && payload[0] == 0x11 && payload[15] == 0x03 && payload[16] == 0x5e {
					// I have no idea what this TTC is but its on top of query
					//simply skiping it
					payload = payload[15:]
				}
				switch payload[0] {
				case 0x03:
					switch payload[1] {
					case 0x5e: //reading query
						query, _ := pascalString(payload[70:])
						context := sql.QueryContext{
							Query:    query,
							Database: o.currentDB,
							User:     o.username,
							Client:   remoteAddrToIP(o.client.RemoteAddr()),
							Time:     time.Now(),
						}
						processContext(context)
					case 0x76: // Reading username
						val, _ := pascalString(payload[19:])
						o.username = val
						logger.Debugf("Username: %s", o.username)
					}
				}
			}
		}

		_, err = o.server.Write(buf)
		if err != nil || eof {
			return err
		}

		err = readWrite(o.server, o.client, o.readPacket)
		if err != nil {
			return err
		}
	}
}

//wrapper around our classic readPacket to handle segmented packets
func (o *Oracle) readPacket(c io.Reader) (buf []byte, err error) {
	buf, err = o.reader(c)
	if err != nil {
		return
	}
	packetLen := int(buf[0])*256 + int(buf[1])
	for {
		if len(buf) == packetLen {
			break
		}
		var b []byte
		b, err = o.reader(c)
		if err != nil {
			return
		}
		buf = append(buf, b...)
	}
	return
}
