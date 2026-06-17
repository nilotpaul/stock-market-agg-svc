package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var (
		dbHost     = os.Getenv("DB_HOST")
		dbKeyspace = os.Getenv("DB_KEYSPACE")
		listenAddr = os.Getenv("API_LISTEN_ADDR")
	)
	if len(dbHost) == 0 {
		log.Fatal("db host environment variable not provided")
	}
	if len(dbKeyspace) == 0 {
		log.Fatal("db keyspace environment variable not provided")
	}
	if len(listenAddr) == 0 {
		log.Fatal("api listen addr environment variable not provided")
	}

	sess, err := OpenDBSession(dbHost, dbKeyspace)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	start, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2026-01-01 09:16:00",
	)
	end, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2026-01-01 09:15:00",
	)

	candles, err := GetCandles(sess, "TCS", start, end)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", candles)
	fmt.Println(len(candles))
}
