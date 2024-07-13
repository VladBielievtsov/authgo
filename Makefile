dev:
	@go run cmd/main.go

run: build
	@./bin/authgo

build:
	@go build -o bin/authgo cmd/main.go

drop:
	@go run cmd/drop/main.go