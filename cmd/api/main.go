package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", conf.Port),
		GracefulTimeout: conf.GracefulTimeout,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
