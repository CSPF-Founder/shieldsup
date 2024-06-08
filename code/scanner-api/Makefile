build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=release -ldflags="-s -w" -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...
	
