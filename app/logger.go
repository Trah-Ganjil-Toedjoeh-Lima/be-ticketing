package app

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func NewLogger(config *config.AppConfig) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	file, _ := os.OpenFile("./storage/logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	output := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(output)
	if config.IsProduction == "false" {
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
