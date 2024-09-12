generate:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/review.proto

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

dev:
	make generate
	goreload main.go