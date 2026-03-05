package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ariefsibuea/flight-aggregator/configs"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	conf := configs.Get()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "OK",
		})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		address := fmt.Sprintf(":%d", conf.Port)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
