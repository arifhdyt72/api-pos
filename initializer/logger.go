package initializer

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logFile := &lumberjack.Logger{
		Filename:   "api.log", // Path to the log file
		MaxSize:    10,        // Max size in megabytes before rotation
		MaxBackups: 5,         // Max number of old log files to retain
		MaxAge:     30,        // Max number of days to retain old log files
		Compress:   true,      // Whether to compress old log files
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(logFile)
	logrus.SetLevel(logrus.InfoLevel)
}
