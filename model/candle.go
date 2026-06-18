package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const CandleDatetimeLayout = "2006-01-02 15:04:05"

type Candle struct {
	Symbol   string    `json:"symbol"`
	DateTime time.Time `json:"datetime"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
}

type GetCandlesRequest struct {
	Symbol    string    `json:"symbol"`
	Timeframe string    `json:"timeframe"`
	Start     time.Time `json:"start_date"`
	End       time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
}

type GetCandlesResponse struct {
	Symbol    string    `json:"symbol"`
	Timeframe string    `json:"timeframe"`
	Candles   []*Candle `json:"candles"`
	Count     int       `json:"count"`
}

func (cr *GetCandlesRequest) ParseAndValidate(r *http.Request) error {
	q := r.URL.Query()

	cr.Symbol = q.Get("symbol")
	if len(cr.Symbol) == 0 {
		return fmt.Errorf("missing symbol")
	}

	cr.Timeframe = q.Get("timeframe")
	if len(cr.Timeframe) == 0 {
		return fmt.Errorf("missing timeframe")
	}

	start, err := time.Parse(CandleDatetimeLayout, q.Get("start_date"))
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}
	cr.Start = start

	end, err := time.Parse(CandleDatetimeLayout, q.Get("end_date"))
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}
	cr.End = end

	if end.Before(start) {
		return fmt.Errorf("end_date must be after start_date")
	}

	limitStr := q.Get("limit")
	if len(limitStr) == 0 {
		cr.Limit = 100
	} else {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("invalid limit")
		}
		if limit <= 0 {
			return fmt.Errorf("limit must be greater than 0")
		}
		if limit > 1000 {
			return fmt.Errorf("limit cannot exceed 1000")
		}
		cr.Limit = limit
	}

	return nil
}
