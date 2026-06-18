package service

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/nilotpaul/stock-market-agg-svc/model"
	"github.com/nilotpaul/stock-market-agg-svc/server/repository"
	"github.com/stretchr/testify/assert"
)

func TestCandleService_aggregateCandles(t *testing.T) {
	candlesData, err := os.ReadFile("testdata/candles.json")
	if err != nil {
		t.Fatal(err)
	}

	expectedCandles, err := os.ReadFile("testdata/expectedCandles.json")
	if err != nil {
		t.Fatal(err)
	}

	var expected []*model.Candle
	if err := json.Unmarshal(expectedCandles, &expected); err != nil {
		t.Fatal(err)
	}

	repo := newTestCandleStore(candlesData)
	svc := NewCandleService(repo)

	got, err := svc.GetCandles(
		t.Context(),
		"TCS",
		mustTime("2026-01-01 09:15:00"),
		mustTime("2026-01-01 09:24:00"),
		"5m",
		100,
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, got)
}

func mustTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

type testCandleStore struct {
	data []byte
}

func newTestCandleStore(b []byte) repository.CandleStorer {
	return &testCandleStore{
		data: b,
	}
}

func (t *testCandleStore) GetCandles(
	ctx context.Context,
	symbol string,
	start time.Time,
	end time.Time,
	limit int,
) ([]*model.Candle, error) {
	var candles []*model.Candle
	if err := json.Unmarshal(t.data, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}
