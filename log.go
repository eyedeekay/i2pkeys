package i2pkeys

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var (
	log  *logrus.Logger
	once sync.Once
)

func InitializeI2PKeysLogger() {
	once.Do(func() {
		log = logrus.New()
		// We do not want to log by default
		log.SetOutput(ioutil.Discard)
		log.SetLevel(logrus.PanicLevel)
		// Check if DEBUG_I2P is set
		if logLevel := os.Getenv("DEBUG_I2P"); logLevel != "" {
			log.SetOutput(os.Stdout)
			switch strings.ToLower(logLevel) {
			case "debug":
				log.SetLevel(logrus.DebugLevel)
			case "warn":
				log.SetLevel(logrus.WarnLevel)
			case "error":
				log.SetLevel(logrus.ErrorLevel)
			default:
				log.SetLevel(logrus.DebugLevel)
			}
			log.WithField("level", log.GetLevel()).Debug("Logging enabled.")
		}
	})
}

// GetI2PKeysLogger returns the initialized logger
func GetI2PKeysLogger() *logrus.Logger {
	if log == nil {
		InitializeI2PKeysLogger()
	}
	return log
}

func init() {
	InitializeI2PKeysLogger()
}
