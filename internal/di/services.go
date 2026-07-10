package di

import (
	"source-base/internal/ent/repository"
	"source-base/internal/service"

	"github.com/google/wire"
)

var RepoServiceSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)
