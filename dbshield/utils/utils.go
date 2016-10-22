package utils

import "net"

//DBMS interface should get implemented with every added DBMS(MySQL, Postgre & etc.) structure
type DBMS interface {
	DefaultPort() uint
	Close()
	SetReader(func(net.Conn) ([]byte, error))
	Handler() error
	SetSockets(net.Conn, net.Conn)
	SetCertificate(string, string) error
}
