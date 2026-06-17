package main

import (
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Candle struct {
	Symbol   string
	DateTime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
}

func OpenDBSession(host, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace

	return cluster.CreateSession()
}

func GetCandles(sess *gocql.Session, symbol string, start, end time.Time) ([]Candle, error) {
	q := sess.Query(`
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
		AND datetime <= ?`,
		symbol, start, end)

	var (
		it = q.Iter()

		candles []Candle
		c       Candle
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
		candles = append(candles, c)
	}

	err := it.Close()
	return candles, err
}
