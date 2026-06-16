APP_NAME=app

runserver:
	@go build -o bin/$(APP_NAME)-server ./server
	@bin/$(APP_NAME)-server

runclient:
	@go build -o bin/$(APP_NAME)-client ./client
	@bin/$(APP_NAME)-client

test:
	@go test -v -count=1 ./...
