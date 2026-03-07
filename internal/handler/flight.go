package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	"github.com/ariefsibuea/flight-aggregator/internal/usecase"

	"github.com/labstack/echo/v4"
)

type flightHandler struct {
	usecase         usecase.FlightUsecase
	providerTimeout time.Duration
}

func InitFlightHandler(e *echo.Group, usecase usecase.FlightUsecase, providerTimeout time.Duration) {
	handler := &flightHandler{
		usecase:         usecase,
		providerTimeout: providerTimeout,
	}

	e.POST("/flights/search", handler.searchFlights)
}

func (f *flightHandler) searchFlights(c echo.Context) error {
	var req model.SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid request body",
		})
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), f.providerTimeout)
	defer cancel()

	res, err := f.usecase.SearchFlights(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "failed to search flights",
		})
	}

	return c.JSON(http.StatusOK, res)
}
