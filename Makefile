.PHONY: vet run_server run_client build_server build_client

vet:
	go vet ./...

run_server:
	cd ./cmd/server && go run .

run_client:
	cd ./cmd/client && go run .

build_server: vet
	go build -o server ./cmd/server

build_client: vet
	go build -o client ./cmd/client

