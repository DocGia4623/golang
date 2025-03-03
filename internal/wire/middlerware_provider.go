package wire

import (
	"testwire/internal/middleware"

	"github.com/google/wire"
)

var MiddlerwareSet = wire.NewSet(
	middleware.NewMiddleware,
)
