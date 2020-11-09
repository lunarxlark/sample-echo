package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(customLogger()))
	e.Use(middleware.Recover())

	// Handler
	e.GET("/greet", greet)
	e.GET("/greet:greeting", greet)
	e.GET("/healthcheck", healthcheck)

	e.Logger.Fatal(e.Start(":80"))
}

var ResponseLogForamt = `{` +
	`"time":"${time_custom}",` +
	`"id":"${id}",` +
	`"remote_ip":"${remote_ip}",` +
	`"host":"${host}",` +
	`"method":"${method}",` +
	`"uri":"${uri}",` +
	`"user_agent":"${user_agent}",` +
	`"status":${status},` +
	`"error":"${error}",` +
	`"latency":${latency},` +
	`"latency_human":"${latency_human}",` +
	`"bytes_in":${bytes_in},` +
	`"bytes_out":${bytes_out},` +
	`"forwarded-for":"${header:x-forwarded-for}",` +
	`"same-as-id":${header:X-Request-Id},` +
	`"query":${query:lang}` +
	`}`

func customLogger() middleware.LoggerConfig {
	cl := middleware.DefaultLoggerConfig
	cl.Skipper = customeSkipper
	cl.Format = ResponseLogForamt
	cl.CustomTimeFormat = "2006/01/02 15:04:05.00000"

	return cl
}

func customeSkipper(c echo.Context) bool {
	if c.Path() == "/healthcheck" {
		return true
	}
	if os.Getenv("ENV") == "auto-test" {
		return true
	} else {
		return false
	}
}

func greet(c echo.Context) error {
	switch c.QueryParam("lang") {
	case "jp":
		return c.String(http.StatusOK, "こんにちは\n")
	case "en":
		return c.String(http.StatusOK, "Hello World\n")
	default:
		return c.String(http.StatusOK, "ジャンボ!!\n")
	}
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "I'm fine\n")
}
