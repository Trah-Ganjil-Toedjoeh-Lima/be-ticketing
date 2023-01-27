package util

import (
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

func (u *SnapUtil) HandleCallback() {

}
