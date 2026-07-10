package di

import (
	"github.com/google/wire"
)

var AppProviderSet = wire.NewSet(
	InfraSet,
	RepoServiceSet,
	WebSet,
)
