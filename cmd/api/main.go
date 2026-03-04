package main

import (
	"fmt"
	"net/http"

	"github.com/ariefsibuea/flight-aggregator/configs"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	conf := configs.Get()

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "OK",
		})
	})

	if err := e.Start(fmt.Sprintf(":%d", conf.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
