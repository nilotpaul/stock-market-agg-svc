package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nilotpaul/stock-market-agg-svc/model"
)

func main() {
	godotenv.Load()
	var (
		symbol = flag.String("symbol", "", "stock symbol")
		tf     = flag.String("timeframe", "1m", "timeframe")
		start  = flag.String("start_date", "", "start datetime")
		end    = flag.String("end_date", "", "end datetime")
	)
	flag.Parse()

	// can also be taken as flag
	listenAddr := os.Getenv("API_LISTEN_ADDR")
	if len(listenAddr) == 0 {
		log.Fatal("api listen addr environment variable not provided")
	}

	// http.NewRequestWithContext should be used
	// but as our API already has ctx timeouts, i'm keeping this simple.
	resp, err := http.Get(buildURL(listenAddr, *symbol, *start, *end, *tf))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		var apiErr model.ApiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			log.Fatalf("failed to parse error response body: %v", err)
		}
		log.Fatalf("failed with status: %d and err: %s", apiErr.Status, apiErr.Message)
	}

	var result model.GetCandlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("failed to parse response body: %v", err)
	}

	printResult(result)
}

func printResult(result model.GetCandlesResponse) {
	fmt.Println("\n=== Fetched Candle Data ===")
	fmt.Printf(
		"\nSymbol: %s | Timeframe: %s | Total Candles: %d\n\n",
		result.Symbol,
		result.Timeframe,
		len(result.Candles),
	)

	for i, c := range result.Candles {
		fmt.Printf(
			"%d %s | O: %.2f | H: %.2f | L: %.2f | C: %.2f | V: %d\n",
			i+1,
			c.DateTime.Format(time.RFC3339),
			c.Open,
			c.High,
			c.Low,
			c.Close,
			c.Volume,
		)
	}
	fmt.Println("\n===========================")
}

func buildURL(listenAddr, symbol, start, end, tf string) string {
	url, _ := url.Parse(fmt.Sprintf("http://127.0.0.1%s/api/v1/candles", listenAddr))
	q := url.Query()

	q.Set("symbol", symbol)
	q.Set("start_date", start)
	q.Set("end_date", end)
	q.Set("timeframe", tf)
	url.RawQuery = q.Encode()

	return url.String()
}
