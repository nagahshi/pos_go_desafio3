createmigration:
	migrate create --ext=mysql -dir=internal/infra/database/migrations -seq
grpcgen:
	protoc --proto_path=internal/infra/grpc/protofiles/ internal/infra/grpc/protofiles/*.proto --go_out=plugins=grpc:.
graphqlgen:
	go run github.com/99designs/gqlgen generate
.PHONY: migrate protoc