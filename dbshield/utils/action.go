package utils

import (
	"net"

	"github.com/nim4/DBShield/dbshield/logger"
)

//ActionDrop will close the connection
func ActionDrop(c net.Conn) error {
	logger.Warning("Dropping connection")
	return nil
}
