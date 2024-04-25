run:
	go run cmd/app/main.go --config=./config/local.yaml

lint:
	golangci-lint run ./...