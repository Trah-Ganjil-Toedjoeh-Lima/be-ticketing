package util

import "github.com/sirupsen/logrus"

type LogUtil struct {
	Log *logrus.Logger
}

func NewLogUtil(log *logrus.Logger) *LogUtil {
	return &LogUtil{Log: log}
}

func (u *LogUtil) BasicLog(err error, occurrence string) {
	u.Log.
		WithField("occurrence", occurrence).
		Error(err.Error())
}

func (u *LogUtil) ControllerResponseLog(err error, occurrence string, clientIp string, clientId uint64) {
	u.Log.
		WithField("occurrence", occurrence).
		WithField("client_ip", clientIp).
		WithField("client_id", clientId).
		Info(err.Error())
}
