BIN := "./bin/rotator"
DB_CONN := "postgres://postgres:postgres@localhost:5432/banners_rotator?sslmode=disable"

build:
	go build -v -o $(BIN) ./cmd/rotator

run: build
	$(BIN) --config ./configs/config.example.json

migrate:
	goose --dir=migrations postgres ${DB_CONN} up