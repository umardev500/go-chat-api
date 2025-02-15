wire:
	wire ./internal/injector

generate:
	protoc --go_out . --go-grpc_out . api/proto/*.proto