package repository

import (
	"context"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/nilotpaul/stock-market-agg-svc/model"
)

type CandleStorer interface {
	GetCandles(ctx context.Context, symbol string, start time.Time, end time.Time, limit int) ([]*model.Candle, error)
}

type CassandraRepo struct {
	session *gocql.Session
}

func NewCassandra(sess *gocql.Session) *CassandraRepo {
	return &CassandraRepo{
		session: sess,
	}
}

func (cr *CassandraRepo) GetCandles(ctx context.Context, symbol string, start, end time.Time, limit int) ([]*model.Candle, error) {
	q := cr.session.Query(`
		SELECT
			symbol,
			datetime,
			open,
			high,
			low,
			close,
			volume
		FROM stock_keyspace.candles
		WHERE symbol = ?
		AND datetime >= ?
		AND datetime <= ?
		LIMIT ?`,
		symbol, start, end, limit)

	var (
		it = q.IterContext(ctx)

		candles []*model.Candle
		c       model.Candle
	)

	for it.Scan(
		&c.Symbol,
		&c.DateTime,
		&c.Open,
		&c.High,
		&c.Low,
		&c.Close,
		&c.Volume,
	) {
		candles = append(candles, &c)
	}

	err := it.Close()
	return candles, err
}
