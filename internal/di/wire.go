//go:build wireinject
// +build wireinject

package di

import (
	"github.com/himbo22/source-base/internal/config"

	"github.com/google/wire"
)

func InitApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		AppProviderSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
