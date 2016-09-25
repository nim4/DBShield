package dbms

import (
	"bytes"
	"net"

	"github.com/nim4/DBShield/dbshield/config"

	"github.com/nim4/DBShield/dbshield/logger"
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

func getMask(key string, data []byte) ([]byte, bool) {
	if mask, ok := config.Config.Masks[key]; ok {
		initLen := len(data)
		logger.Debug("Masking")
		data = mask.MatchExp.ReplaceAll(data, mask.ReplaceExp)
		if len(data) < initLen {
			data = append(data, bytes.Repeat(mask.PaddingCharacter, initLen-len(data))...)
		}
		return data[:initLen], true
	}
	return nil, false
}
