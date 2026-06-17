package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/joho/godotenv"
)

const stockDataCSVPath = "data/stock_data.csv"

const (
	FieldSymbol = iota
	FieldDatetime
	FieldOpen
	FieldHigh
	FieldLow
	FieldClose
	FieldVolume
)

func main() {
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		log.Fatal("db host environment variable not provided")
	}

	file, err := os.Open(stockDataCSVPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cluster := gocql.NewCluster(dbHost)

	sess, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("creating keyspace and table")
	if err := setupDB(sess); err != nil {
		log.Fatal(err)
	}

	fmt.Println("seeding db with csv data")
	if err := seedDB(file, sess); err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}

func seedDB(fr io.Reader, sess *gocql.Session) error {
	const query = `
		INSERT INTO stock_keyspace.candles (
			symbol,
			datetime,
			open,
			high,
			low,
			close,
			volume
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	r := csv.NewReader(fr)
	_, _ = r.Read() // skip the headers

	for {
		record, err := r.Read()
		if err != nil {
			// exit and return when EOF
			if err == io.EOF {
				return nil
			}
			return err
		}

		// parse time with csv data layout
		dt, err := time.Parse("2006-01-02 15:04:05", record[FieldDatetime])
		if err != nil {
			return fmt.Errorf("invalid datetime %q: %w", record[FieldDatetime], err)
		}
		// convert from str to float
		open, err := strconv.ParseFloat(record[FieldOpen], 64)
		if err != nil {
			return err
		}
		high, err := strconv.ParseFloat(record[FieldHigh], 64)
		if err != nil {
			return err
		}
		low, err := strconv.ParseFloat(record[FieldLow], 64)
		if err != nil {
			return err
		}
		closePrice, err := strconv.ParseFloat(record[FieldClose], 64)
		if err != nil {
			return err
		}
		volume, err := strconv.ParseInt(record[FieldVolume], 10, 64)
		if err != nil {
			return err
		}

		err = sess.Query(
			query,
			record[FieldSymbol],
			dt,
			open,
			high,
			low,
			closePrice,
			volume,
		).Exec()
		if err != nil {
			return err
		}
	}
}

// setupDB will create keyspace and table.
func setupDB(sess *gocql.Session) error {
	queries := []string{
		`
		CREATE KEYSPACE IF NOT EXISTS stock_keyspace
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
		`,
		`
		CREATE TABLE IF NOT EXISTS stock_keyspace.candles (
			symbol text,
			datetime timestamp,
			open double,
			high double,
			low double,
			close double,
			volume bigint,
			PRIMARY KEY ((symbol), datetime)
		) WITH CLUSTERING ORDER BY (datetime ASC)
		`,
	}

	for _, q := range queries {
		if err := sess.Query(q).Exec(); err != nil {
			return err
		}
	}

	return nil
}
