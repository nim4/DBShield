package dbms

import (
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
		go training.AddToTrainingSet(context)
	} else {
		if config.Config.ActionFunc != nil && !training.CheckQuery(context) {
			return config.Config.ActionFunc()
		}
	}
	return nil
}
