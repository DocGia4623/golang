package wire

import (
	"testwire/internal/controller"

	"github.com/google/wire"
)

// ControllerSet chứa tất cả controller
var ControllerSet = wire.NewSet(
	controller.NewAuthenticationController,
	controller.NewUserController,
)
