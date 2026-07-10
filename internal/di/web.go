package di

import (
	"source-base/internal/bootstrap"
	http "source-base/internal/controller/http"

	"github.com/google/wire"
)

var WebSet = wire.NewSet(
	http.NewUserController,
	bootstrap.InitRouter,
	bootstrap.InitEcho,
)
