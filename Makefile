APP_NAME=app

runserver:
	@go build -o bin/$(APP_NAME)-server ./server
	@bin/$(APP_NAME)-server

runclient:
	@go build -o bin/$(APP_NAME)-client ./client
	@bin/$(APP_NAME)-client

init:
	@go run ./server/script/main.go

cli:
	go run ./client/main.go \
    -symbol=TCS \
    -start_date="2026-01-01 09:16:00" \
    -end_date="2026-01-01 09:21:00" \
    -timeframe=1m

test:
	@go test -v -count=1 ./...
