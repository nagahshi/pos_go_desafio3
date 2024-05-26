createmigration:
	migrate create --ext=mysql -dir=internal/infra/database/migrations -seq
grpcgen:
	protoc --proto_path=internal/infra/grpc/protofiles/ internal/infra/grpc/protofiles/*.proto --go_out=plugins=grpc:.
graphqlgen:
	go run github.com/99designs/gqlgen generate
run:
	@go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
build:
	go build cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go -o main && ./wait-for-it.sh localhost:8080

.PHONY: migrate protoc