package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ariefsibuea/flight-aggregator/internal/model"
	pkgerr "github.com/ariefsibuea/flight-aggregator/internal/pkg/errors"
	"github.com/ariefsibuea/flight-aggregator/internal/usecase"

	"github.com/labstack/echo/v4"
)

type flightHandler struct {
	usecase         usecase.FlightUsecase
	providerTimeout time.Duration
}

func InitFlightHandler(e *echo.Group, uc usecase.FlightUsecase, providerTimeout time.Duration) {
	handler := &flightHandler{
		usecase:         uc,
		providerTimeout: providerTimeout,
	}

	e.POST("/flights/search", handler.searchFlights)
}

func (f *flightHandler) searchFlights(c echo.Context) error {
	var req model.SearchRequest
	if err := c.Bind(&req); err != nil {
		return pkgerr.BadRequestError("invalid request body")
	}

	if err := req.Validate(); err != nil {
		return pkgerr.ValidationErrorf("validation failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), f.providerTimeout)
	defer cancel()

	res, err := f.usecase.SearchFlights(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return pkgerr.BadRequestError("search request timed out")
		}
		return err
	}

	return Success(c, http.StatusOK, res)
}
