//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/config"
	"github.com/google/wire"
)

func InitializeServer() *app.Server {
	wire.Build(
		config.NewAppConfig,
		app.NewServer,
	)
	return nil
}
