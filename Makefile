generate:
	protoc --go_out=. --go_out=paths=source_relative  --go-grpc_out=. --go-grpc_out=paths=source_relative proto/*.proto
run:
	go run cmd/main.go