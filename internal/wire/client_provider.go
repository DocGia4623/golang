package wire

import (
	"testwire/config"

	"github.com/google/wire"
)

var ClientSet = wire.NewSet(
	config.NewElasticClient,
)
