build:
# Compile the Go application, this runs the main file
	@go build -o bin/personal-database ./cmd/main.go

run: 
# Run the Go application
	@./bin/personal-database

test:
# Run tests
	@go test -v ./...