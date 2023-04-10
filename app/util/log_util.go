package util

import "github.com/sirupsen/logrus"

type LogUtil struct {
	Logrus *logrus.Logger
}

func NewLogUtil(log *logrus.Logger) *LogUtil {
	return &LogUtil{Logrus: log}
}

// BasicLog create basic error log using logrus
func (u *LogUtil) BasicLog(err error, occurrence string) {
	u.Logrus.
		WithField("occurrence", occurrence).
		Error(err.Error())
}

// ControllerResponseLog create basic error (specially for controller) log using logrus
func (u *LogUtil) ControllerResponseLog(err error, occurrence string, clientIp string, clientId uint64) {
	u.Logrus.
		WithField("occurrence", occurrence).
		WithField("client_ip", clientIp).
		WithField("client_id", clientId).
		Info(err.Error())
}
