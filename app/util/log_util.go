package util

import "github.com/sirupsen/logrus"

type LogUtil struct {
	log *logrus.Logger
}

func NewLogUtil(log *logrus.Logger) *LogUtil {
	return &LogUtil{log: log}
}

func (u *LogUtil) BasicLog(err error, occurrence string) {
	u.log.
		WithField("occurrence", occurrence).
		Error(err.Error())
}

func (u *LogUtil) ControllerResponseLog(err error, occurrence string, clientIp string, clientId uint64) {
	u.log.
		WithField("occurrence", occurrence).
		Info(err.Error())
}
