package di

import (
	"github.com/himbo22/source-base/internal/bootstrap"
	http "github.com/himbo22/source-base/internal/controller/http"

	"github.com/google/wire"
)

var WebSet = wire.NewSet(
	http.NewUserController,
	bootstrap.InitRouter,
	bootstrap.InitEcho,
)
