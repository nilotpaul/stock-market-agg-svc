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
