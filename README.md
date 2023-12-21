# Johnjud-gateway

## Stack

-   golang
-   gRPC
-   go-fiber

## Getting Started

### Prerequisites

-   golang 1.21 or [later](https://go.dev)
-   docker
-   makefile

### Installation

1. Clone this repo
2. Copy `config.example.yaml` in `config` and paste it in the same directory with `.example` removed from its name.

3. Run `go mod download` to download all the dependencies.

### Running
1. Run `docker-compose up -d`
2. Run `make server` or `go run ./src/.`

### Testing
1. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`
