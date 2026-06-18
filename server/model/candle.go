package model

import (
	"fmt"
	"net/http"
	"time"
)

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
	Symbol    string
	Timeframe string
	Start     time.Time
	End       time.Time
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

	start, err := time.Parse("2006-01-02 15:04:05", q.Get("start_date"))
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}
	cr.Start = start

	end, err := time.Parse("2006-01-02 15:04:05", q.Get("end_date"))
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}
	cr.End = end

	if end.Before(start) {
		return fmt.Errorf("end_date must be after start_date")
	}

	return nil
}
