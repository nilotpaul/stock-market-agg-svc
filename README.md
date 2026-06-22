# Stock Market Data Aggregation Service

## Requirements

- Go
- Apache Cassandra
- Docker (used docker & compose file is in root)

## Running the Project

> **Note:** All commands below are for Linux.
>
> All required commands are available in the Makefile. If `make` is not installed, use the equivalent Go commands shown below.

### 1. Configure Environment

If using docker:

In project root run:
```bash
docker compose up -d
```

Update `.env`:

```env
# NOTE: Change the IP or use localhost
DB_HOST=192.168.0.129:9042
```

### 2. Initialize Database

Creates the keyspace, table, and seeds sample data.

```bash
make init
```

or

```bash
go run ./server/script/main.go
```

> Initial setup and seeding may take up to 10 seconds.

### 3. Run the API Server

```bash
make runserver
```

or

```bash
go run ./server
```

Visit http://localhost:3000, this will serve both json api and web vue.js client. 

### 4. Test Using the CLI Client

```bash
make cli
```

CLI options can be modified in the Makefile.

Or run manually:

```bash
go run ./client/main.go \
    -symbol=TCS \
    -start_date="2026-01-01 09:16:00" \
    -end_date="2026-01-01 09:21:00" \
    -timeframe=1m
```

### 5. Test Using cURL

```bash
curl -X GET -G \
  -H "Accept: application/json" \
  --data-urlencode "symbol=TCS" \
  --data-urlencode "start_date=2026-01-01 09:16:00" \
  --data-urlencode "end_date=2026-01-01 09:21:00" \
  --data-urlencode "timeframe=1m" \
  "http://localhost:3000/api/v1/candles"
```

### 6. Run Tests

```bash
make test
```

or

```bash
go test -v -count=1 ./...
```

## API Documentation

### Health Check

#### Request

```http
GET /api/v1/health
```

#### Example

```bash
curl -X GET http://localhost:3000/api/v1/health
```

#### Response

```json
"OK"
```

---

### Get Candles

Returns candle data for a symbol within a date range. Data can optionally be aggregated into larger timeframes.

#### Request

```http
GET /api/v1/candles
```

#### Query Parameters

| Parameter  | Type   | Required | Description                                             |
| ---------- | ------ | -------- | ------------------------------------------------------- |
| symbol     | string | Yes      | Stock symbol                                            |
| timeframe  | string | Yes      | Candle timeframe (`1m`, `5m`, `15m`, `30m`, `1h`, `1d`) |
| start_date | string | Yes      | Start datetime (`YYYY-MM-DD HH:mm:ss`)                  |
| end_date   | string | Yes      | End datetime (`YYYY-MM-DD HH:mm:ss`)                    |
| limit      | int    | No       | Maximum number of records returned (default: 100)       |

#### Example

```bash
curl -X GET -G \
  --data-urlencode "symbol=TCS" \
  --data-urlencode "start_date=2026-01-01 09:15:00" \
  --data-urlencode "end_date=2026-01-01 09:30:00" \
  --data-urlencode "timeframe=5m" \
  "http://localhost:3000/api/v1/candles"
```

#### Example Response

```json
{
  "symbol": "TCS",
  "timeframe": "5m",
  "candles": [
    {
      "symbol": "TCS",
      "datetime": "2026-01-01T09:15:00Z",
      "open": 3215,
      "high": 3225.3,
      "low": 3208.3,
      "close": 3223.3,
      "volume": 156747
    }
  ],
  "count": 1
}
```

#### Error Response

Can be of status 400, 404, 500.

```json
{
  "status": 404,
  "message": "candles data not found",
  "data": null,
}
```

## Assumptions & Notes

- I’ve kept the implementation simple and mainly focused on the core requirements to keep things simpler for both me and the evaluator.
- The `model` package was kept at the project root because I needed the types for both the server and CLI client. Otherwise, I would've kept it inside the server package.
- For timeframe aggregation, I assumed I needed to convert the requested timeframe into minutes since the provided data was 1m.
- I've skipped graceful shutdown handling for simplicity.
- For pagination, I only implemented a `limit` query parameter for simplicity and time constraints. I understand that Cassandra's preferred pagination approach is to send the previous iterator paging state and use that to fetch the next batch of results.
- For the HTTP layer, I intentionally did not use any libraries for simplicity and used a centralized error handling strategy, which is how I personally like to structure.
- I couldn't implement the caching due to time constraints. My intended approach would be to create a `Cacher` interface and use:

  ```
  symbol + start_time + end_time + timeframe
  ```

  as the cache key and:

  ```go
  []*Candle
  ```

  as the cached value in an in memory store:

  ```go
  map[string][]*Candle
  ```
- In the CLI client, I did not use a context when calling the HTTP service because the API Service already has timeouts. For production, I would use `http.NewRequestWithContext`.
- For logging and observability, I kept things basic and limited to request logging, error logging, and request duration measurements.
