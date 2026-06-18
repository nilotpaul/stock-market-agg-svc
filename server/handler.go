package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nilotpaul/stock-market-agg-svc/model"
	"github.com/nilotpaul/stock-market-agg-svc/server/service"
)

type CandleHandler struct {
	svc service.CandleService
}

func NewCandleHandler(svc service.CandleService) *CandleHandler {
	return &CandleHandler{
		svc: svc,
	}
}

func (h *CandleHandler) HandleGetCandles(w http.ResponseWriter, r *http.Request) error {
	req := new(model.GetCandlesRequest)

	if err := req.ParseAndValidate(r); err != nil {
		return NewBadRequestError(err.Error(), req)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*300)
	defer cancel()

	candles, err := h.svc.GetCandles(
		ctx,
		req.Symbol,
		req.Start, req.End,
		req.Timeframe,
		req.Limit,
	)
	if err != nil {
		if errors.Is(err, service.ErrUnsupportedTimeframe) {
			return NewBadRequestError(err.Error(), req)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return NewInternalServerError("request took too long to respond", nil)
		}
		return err
	}

	if len(candles) == 0 {
		return NewNotFoundError("candles data not found", nil)
	}

	return writeJSON(w, http.StatusOK, model.GetCandlesResponse{
		Symbol:    req.Symbol,
		Timeframe: req.Timeframe,
		Candles:   candles,
		Count:     len(candles),
	})
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, "OK")
}
