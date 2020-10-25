package socket

import (
	"battleship/config"
	"battleship/dto"
	"battleship/middlewares"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v2"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Printf("%+v", r)
		return true
	},
} // use default options

func StartSocketServer() {
	e := echo.New()
	e.HideBanner = true
	setHttpMiddlewares(e)
	e.Logger = lecho.From(log.Logger)
	e.GET("/", func(c echo.Context) error {
		return connect(c.Response(), c.Request())
	})

	httpConfig := &http.Server{
		Addr: fmt.Sprintf(":%s", config.C.SocketPort),
	}

	log.Fatal().Err(e.StartServer(httpConfig)).Msg("")
}

func setHttpMiddlewares(e *echo.Echo) {
	e.Use(middleware.BodyDump(middlewares.BodyDumper))
	e.Use(middlewares.LogMiddleware())
}

func connect(w http.ResponseWriter, r *http.Request) error {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Msg("error in upgrading:" + err.Error())
		return err
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Warn().Msg("error in reading message:" + err.Error())
			break
		}
		log.Debug().Msg("recv: " + string(message))
		event := new(dto.Event)
		err = json.Unmarshal(message, event)
		if err != nil {
			log.Error().Str("message", string(message)).Msg("message format is not Event")
			//_ = c.Close()
		}

		err = handleEvents(*event, c)
		if err != nil {
			//_ = c.Close()
		}

		//err = c.WriteMessage(1, message)
		//if err != nil {
		//	log.Warn().Msg("write:" + err.Error())
		//	break
		//}
	}
	return nil
}
