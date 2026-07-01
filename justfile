set dotenv-load

migrate-up:
    migrate -path ./migrations -database "$DB_URL" up

migrate-down:
    migrate -path ./migrations -database "$DB_URL" down

run:
    go run main.go