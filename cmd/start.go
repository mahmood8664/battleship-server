package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"battleship/config"
	"battleship/di"
)

var startCMD = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var (
	checkHealth = di.CreateCheckHealthController()
)

func startServer() {
	e := echo.New()
	e.HideBanner = true
	setMiddlewares(e)
	setEndpoints(e)
	e.Debug = config.C.Logging.Level == "debug"
	httpConfig := &http.Server{
		Addr: fmt.Sprintf(":%s", config.C.Port),
	}
	logrus.Fatal(e.StartServer(httpConfig))
}

func setMiddlewares(e *echo.Echo) {

}

func setEndpoints(e *echo.Echo) {
	e.GET("/", checkHealth.CheckHealth())
}
