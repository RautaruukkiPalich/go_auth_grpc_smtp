run:
	go run cmd/mail/main.go --config=./config/local.yaml

lint:
	golangci-lint run ./...