//+build wireinject

package di

import (
	"github.com/google/wire"
	"battleship/controllers"
)

func CreateCheckHealthController() *controllers.CheckHealthController {
	panic(wire.Build(
		controllers.NewCheckHealthController,
	))
}
