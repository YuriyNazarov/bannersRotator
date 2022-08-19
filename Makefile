BIN := "./bin/rotator"
DB_CONN := "postgres://postgres:postgres@localhost:5432/banners_rotator?sslmode=disable"

build:
	go build -v -o $(BIN) ./cmd/rotator

run: build
	$(BIN)

migrate:
	goose --dir=migrations postgres ${DB_CONN} up

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race ./internal/... -count 100