package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nilotpaul/stock-market-agg-svc/model"
	"github.com/nilotpaul/stock-market-agg-svc/server/repository"
)

var (
	ErrUnsupportedTimeframe = errors.New("unsupported timeframe")
)

type CandleService interface {
	GetCandles(ctx context.Context, symbol string, start, end time.Time, tf string, limit int) ([]*model.Candle, error)
}

type candleService struct {
	repo repository.CandleStorer
}

func NewCandleService(repo repository.CandleStorer) CandleService {
	return &candleService{
		repo: repo,
	}
}

func (svc *candleService) GetCandles(ctx context.Context, symbol string, start, end time.Time, tf string, limit int) ([]*model.Candle, error) {
	candles, err := svc.repo.GetCandles(ctx, symbol, start, end, limit)
	if err != nil {
		return nil, err
	}

	return svc.aggregateCandles(candles, tf)
}

func (svc *candleService) aggregateCandles(candles []*model.Candle, tf string) ([]*model.Candle, error) {
	mins, err := svc.timeframeToMinutes(tf)
	if err != nil {
		return nil, err
	}
	if mins == 1 {
		return candles, nil
	}

	var result []*model.Candle
	for i := 0; i < len(candles); i += mins {
		end := i + mins
		if end > len(candles) {
			end = len(candles)
		}

		window := candles[i:end]
		agg := &model.Candle{
			Symbol:   window[0].Symbol,
			DateTime: window[0].DateTime,
			Open:     window[0].Open,
			High:     window[0].High,
			Low:      window[0].Low,
			Close:    window[len(window)-1].Close,
		}

		for _, c := range window {
			if c.High > agg.High {
				agg.High = c.High
			}
			if c.Low < agg.Low {
				agg.Low = c.Low
			}
			agg.Volume += c.Volume
		}

		result = append(result, agg)
	}

	return result, nil
}

func (*candleService) timeframeToMinutes(tf string) (int, error) {
	switch tf {
	case "1m":
		return 1, nil
	case "5m":
		return 5, nil
	case "15m":
		return 15, nil
	case "30m":
		return 30, nil
	case "1h":
		return 60, nil
	case "1d":
		return 24 * 60, nil
	default:
		return 0, fmt.Errorf("%w %s", ErrUnsupportedTimeframe, tf)
	}
}
