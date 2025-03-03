package wire

import (
	"testwire/internal/services"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	services.NewAuthenticationServiceImpl,
	services.NewRefreshTokenServiceImpl,
	services.NewUserServiceImpl,
)
