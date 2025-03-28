package pkg

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Logger *logrus.Logger
var GORMLogger logger.Interface

func InitLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.InfoLevel)
	GORMLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		},
	)
}
