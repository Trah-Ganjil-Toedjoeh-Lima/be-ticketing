package util

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type SnapUtil struct {
	app        *config.AppConfig
	snapClient snap.Client
}

func NewSnapUtil(app *config.AppConfig, logrus *logrus.Logger) *SnapUtil {
	var snapClient snap.Client
	if app.MidtransIsProduction == false {
		snapClient.New(app.ServerKeySandbox, midtrans.Sandbox)
	} else {
		snapClient.New(app.ServerKeyProduction, midtrans.Production)
	}

	if app.IsProduction == true || app.MidtransIsProduction == true {
		snapClient.HttpClient = &midtrans.HttpClientImplementation{
			HttpClient: midtrans.DefaultGoHttpClient,
			Logger:     &SnapLogger{LogLevel: midtrans.LogError, Logrus: logrus},
		}
	} else {
		snapClient.HttpClient = &midtrans.HttpClientImplementation{
			HttpClient: midtrans.DefaultGoHttpClient,
			Logger:     &SnapLogger{LogLevel: midtrans.LogDebug, Logrus: logrus},
		}
	}

	return &SnapUtil{
		app:        app,
		snapClient: snapClient,
	}
}

func (u *SnapUtil) CreateTransaction(request *snap.Request) (*snap.Response, *midtrans.Error) {
	resp, err := u.snapClient.CreateTransaction(request)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (u *SnapUtil) CheckSignature(message map[string]interface{}) error {
	orderId := message["order_id"].(string)
	statusCode := message["status_code"].(string)
	grossAmt := message["gross_amount"].(string)
	signatureKey := message["signature_key"].(string)
	var serverKey string
	if u.app.MidtransIsProduction == false {
		serverKey = u.app.ServerKeySandbox
	} else {
		serverKey = u.app.ServerKeyProduction
	}
	payload := orderId + statusCode + grossAmt + serverKey
	hasher := sha512.New()
	hasher.Write([]byte(payload))
	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))
	if hashStr != signatureKey {
		return errors.New("SIGNATURE KEY NOT MATCH. Signature key: " + signatureKey + " given: " + hashStr)
	}
	return nil
}

type SnapLogger struct {
	LogLevel midtrans.LogLevel
	Logrus   *logrus.Logger
}

// Error : Logs a warning message using Printf conventions.
func (l *SnapLogger) Error(format string, val ...interface{}) {
	if l.LogLevel >= midtrans.LogError {
		l.Logrus.
			WithField("occurrence", "MIDTRANS ERROR").
			WithField("FORMAT", format).
			Error(val)
	}
}

// Info : Logs information message using Printf conventions.
func (l *SnapLogger) Info(format string, val ...interface{}) {
	if l.LogLevel >= midtrans.LogInfo {
		l.Logrus.
			WithField("occurrence", "MIDTRANS INFO").
			WithField("FORMAT", format).
			Info(val)
	}
}

// Debug : Log debug message using Printf conventions.
func (l *SnapLogger) Debug(format string, val ...interface{}) {
	l.Logrus.
		WithField("occurrence", "MIDTRANS DEBUG").
		WithField("FORMAT", format).
		Debug(val)
}
