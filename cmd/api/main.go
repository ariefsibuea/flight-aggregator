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
	"github.com/ariefsibuea/flight-aggregator/internal/cache"
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

	e.HTTPErrorHandler = handler.ErrorHandler

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	})

	apiGroup := e.Group("/api/v1")

	rawProviders := []provider.FlightFetcher{
		garuda.NewClient(),
		lionair.NewClient(),
		batikair.NewClient(),
		airasia.NewClient(),
	}

	providers := make([]provider.FlightFetcher, 0, len(rawProviders))
	for _, p := range rawProviders {
		providers = append(providers,
			provider.NewFlightFetcherLimiter(
				provider.NewFlightFetcherRetrier(p, conf.ProviderMaxRetries, conf.ProviderRetryDelay),
				conf.ProviderRateLimitRPS,
			),
		)
	}

	redisAddr := fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort)
	redisCache := cache.NewRedisCache(redisAddr)

	flightUsecase := usecase.NewFlightUsecase(providers, redisCache, conf.DefaultCacheTTL)
	handler.InitFlightHandler(apiGroup, flightUsecase, conf.ProviderTimeout)

	e.Server.ReadTimeout = conf.ServerReadTimeout
	e.Server.WriteTimeout = conf.ServerWriteTimeout
	e.Server.IdleTimeout = conf.ServerIdleTimeout

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		address := fmt.Sprintf(":%d", conf.ServerPort)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ServerShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
