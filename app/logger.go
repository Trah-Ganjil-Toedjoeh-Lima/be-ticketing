package app

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

func NewLogger(config *config.AppConfig) *logrus.Logger {
	logger := logrus.New()
	output := io.MultiWriter(os.Stdout)
	logger.SetOutput(output)

	if config.IsProduction == false {
		logger.SetReportCaller(true) //for detailing the error's line of code
		logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				return frame.Function, fileName
			},
		})
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
