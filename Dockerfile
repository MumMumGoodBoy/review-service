# Start with a Go base image
FROM golang:1.23.2-bullseye AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install protoc and the Go plugin for protoc
RUN apt update && \
    apt install -y unzip && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v28.3/protoc-28.3-linux-aarch_64.zip && \
    unzip protoc-28.3-linux-aarch_64.zip -d /usr/local bin/protoc && \
    unzip protoc-28.3-linux-aarch_64.zip -d /usr/local 'include/*' && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download


# Copy the rest of the application code
COPY . .

RUN make generate

RUN go build -o myapp .

FROM debian:bullseye-slim

COPY --from=builder /app/myapp /myapp

EXPOSE 3000

CMD ["/myapp"]
