package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ariefsibuea/flight-aggregator/config"
	"github.com/ariefsibuea/flight-aggregator/internal/handler"
	"github.com/ariefsibuea/flight-aggregator/internal/provider"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/airasia"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/batikair"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/garuda"
	"github.com/ariefsibuea/flight-aggregator/internal/provider/lionair"
	"github.com/ariefsibuea/flight-aggregator/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	conf := config.Get()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "OK",
		})
	})

	apiGroup := e.Group("/api/v1")

	providers := []provider.FlightFetcher{
		garuda.NewClient(),
		lionair.NewClient(),
		batikair.NewClient(),
		airasia.NewClient(),
	}
	flightUsecase := usecase.NewFlightUsecase(providers)
	handler.InitFlightHandler(apiGroup, flightUsecase, conf.ProviderTimeout)

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
