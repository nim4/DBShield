package utils

import (
	"io"
	"net"
	"reflect"
)

//DBMS interface should get implemented with every added DBMS(MySQL, Postgre & etc.) structure
type DBMS interface {
	DefaultPort() uint
	Close()
	SetReader(func(io.Reader) ([]byte, error))
	Handler() error
	SetSockets(net.Conn, net.Conn)
	SetCertificate(string, string) error
}

//GenerateDBMS instantiate a new instance of DBMS
func GenerateDBMS(original DBMS) DBMS {
	val := reflect.ValueOf(original)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	newThing := reflect.New(val.Type()).Interface().(DBMS)
	return newThing
}
