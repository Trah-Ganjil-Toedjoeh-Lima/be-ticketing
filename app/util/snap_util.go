package util

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type SnapUtil struct {
	app        *config.AppConfig
	snapClient snap.Client
}

func NewSnapUtil(app *config.AppConfig) *SnapUtil {
	var snapClient snap.Client
	if app.MidtransIsProduction == false {
		snapClient.New(app.ServerKeySandbox, midtrans.Sandbox)
	} else {
		snapClient.New(app.ServerKeySandbox, midtrans.Production)
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

func (u *SnapUtil) HandleCallback() {

}
