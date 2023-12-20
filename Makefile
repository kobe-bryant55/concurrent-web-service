TEST_OUTPUT := coverage.out

# Start db container
start-db:
	docker compose up -d postgres-db

# Stop db container
down-db:
	docker compose stop postgres-db

# Remove db container
rm-db:
	docker compose rm -f postgres-db


# Run all Tests
test:
	go test ./... -v

# Cover tests
cover:
	go test -coverpkg=./... -coverprofile=$(TEST_OUTPUT) ./... ./internal/database
	go tool cover -html=$(TEST_OUTPUT)

# Swagger init
swag:
	swag init -parseDependency -g ./cmd/api/main.go