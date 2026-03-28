# Makefile for godistributedkv

.PHONY: build run test tidy clean protoc

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./... -v

tidy:
	go mod tidy

clean:
	rm -rf bin/ data/

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/kv.proto
