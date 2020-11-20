package middlewares

import (
	"battleship/dto"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

func LogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if strings.HasPrefix(c.Request().RequestURI, "/api/v") {
				req := c.Request()
				res := c.Response()

				start := time.Now()
				if err = next(c); err != nil {
					dto.CustomHTTPErrorHandler(err, c)
				}
				stop := time.Now()

				bytesIn := req.Header.Get(echo.HeaderContentLength)
				if bytesIn == "" {
					bytesIn = "0"
				}

				var logg *zerolog.Event

				switch {
				case 200 <= res.Status && res.Status < 300:
					logg = log.Debug()
				case 300 <= res.Status && res.Status < 400, 0 < res.Status && res.Status < 200:
					logg = log.Info()
				case 500 <= res.Status:
					logg = log.Error()
				default:
					logg = log.Info()
				}

				event := logg.
					Str("request_time", start.Format(time.RFC3339Nano)).
					Str("remote_ip", req.RemoteAddr).
					Str("real_ip", c.RealIP()).
					Str("host", req.Host).
					Str("method", req.Method).
					Str("uri", req.RequestURI).
					Str("user_agent", req.UserAgent()).
					Int("status", res.Status).
					Str("latency", strconv.FormatInt(int64(stop.Sub(start)), 10)).
					Str("latency_human", stop.Sub(start).String()).
					Str("bytes_in", bytesIn).
					Str("bytes_out", strconv.FormatInt(res.Size, 10))

				var errorMessage string
				if err != nil {
					if he, ok := err.(*echo.HTTPError); ok {
						if hs, ok := he.Message.(echo.Map); ok {
							errorMessage = hs["message"].(string)
						} else {
							errorMessage = he.Message.(string)
						}
					} else {
						errorMessage = err.Error()
					}
				}
				if errorMessage != "" {
					event.Str("error_message", errorMessage)
				}
				event.Msg("Http Request")
			} else {
				_ = next(c)
			}
			return
		}
	}
}

func BodyDumper(c echo.Context, reqBody, resBody []byte) {
	if strings.HasPrefix(c.Request().RequestURI, "/api/v") {
		log.Debug().
			Str("uri", c.Request().RequestURI).
			Str("request_body", string(reqBody)).
			Str("response_body", string(resBody)).
			Msg("Http Request")
	}
	return
}
