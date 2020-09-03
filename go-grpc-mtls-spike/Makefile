all: lint build

build: proto_gen
	go build --mod=vendor -o bin/proxy ./cmd/proxy/main.go
	go build --mod=vendor -o bin/greeter ./cmd/greeter/main.go

lint:
	golangci-lint run ./...

vendor:
	go mod tidy
	go mod vendor

proto_gen:
	protoc \
		--go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
		--go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		internal/protobuf/*.proto

.PHONY: all build vendor