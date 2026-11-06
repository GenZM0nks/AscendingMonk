.PHONY: run swagger

run:
	DB_PATH="../local_monk_db/ascendingmonk.db" go run ./cmd/main

swagger:
	swag init -g cmd/main/main.go -o docs/