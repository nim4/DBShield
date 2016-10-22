package dbms

import (
	"crypto/tls"
	"net"

	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

func pascalString(data []byte) (str string, size uint) {
	size = uint(data[0])
	str = string(data[1 : size+1])
	return
}

func remoteAddrToIP(addr net.Addr) string {
	return addr.(*net.TCPAddr).IP.String()
}

func handlePanic() {
	if r := recover(); r != nil {
		logger.Warningf("%v", r)
	}
}

func threeByteBigEndianToInt(data []byte) uint {
	return uint(data[2])*65536 + uint(data[1])*256 + uint(data[0])
}

//processContext will handle context depending on running mode
func processContext(context sql.QueryContext) (err error) {
	logger.Infof("Query: %s", context.Query)

	if config.Config.Learning {
		return training.AddToTrainingSet(context)
	}
	if config.Config.ActionFunc != nil && !training.CheckQuery(context) {
		return config.Config.ActionFunc()
	}
	return nil
}

func turnSSL(client net.Conn, server net.Conn, certificate tls.Certificate) (net.Conn, net.Conn, error) {
	logger.Info("SSL connection")
	tlsConnClient := tls.Server(client, &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: true,
	})
	if err := tlsConnClient.Handshake(); err != nil {
		return nil, nil, err
	}
	logger.Debug("Client handshake done")

	tlsConnServer := tls.Client(server, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err := tlsConnServer.Handshake(); err != nil {
		return nil, nil, err
	}
	logger.Debug("Server handshake done")
	return tlsConnClient, tlsConnServer, nil
}

func readWrite(src, dst net.Conn, reader func(net.Conn) ([]byte, error)) error {
	//Read result from server
	buf, err := reader(src)
	if err != nil {
		return err
	}

	//Send result to client
	_, err = dst.Write(buf)
	return err
}
