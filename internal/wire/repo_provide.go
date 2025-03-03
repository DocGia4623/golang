package wire

import (
	"testwire/internal/repository"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	repository.NewUserRepositoryImpl,
	repository.NewRefreshTokenRepositoryImpl,
	repository.NewPermissionRepositoryImpl,
)
