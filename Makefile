TEST_OUTPUT := coverage.out

# Start db container
db-start:
	docker compose up -d postgres-db

# Stop db container
db-down:
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

# Install linter dependencies
lint-dep:
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/bombsimon/wsl/v4/cmd...@master

# Run the linter on the specified path
lint:
	go mod tidy
	go vet ./...
	go fmt ./...
	gci write -s standard -s default -s "prefix(github.com/mehmettalhaseker/concurrent-web-service)" .
	gofumpt -l -w .
	wsl -fix ./...
	golangci-lint run $(p)
.PHONY: lint