.PHONY: run swagger

run:
	DB_PATH="whoknows.db" go run ./cmd/main

swagger:
	swag init -g cmd/main/main.go -o docs/