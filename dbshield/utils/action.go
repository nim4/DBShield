package utils

import (
	"net"

	"../logger"
)

//ActionDrop will close the connection
func ActionDrop(c net.Conn) error {
	logger.Warning("Dropping connection")
	return nil
}
