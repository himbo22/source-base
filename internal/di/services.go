package di

import (
	"github.com/himbo22/source-base/internal/ent/repository"
	"github.com/himbo22/source-base/internal/service"

	"github.com/google/wire"
)

var RepoServiceSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)
