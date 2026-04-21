package app

import (
	"go.uber.org/fx"
)

var Kernel = fx.Module(
	"bootstrap",
	fx.Provide(),
)
