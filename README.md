# Johnjud-gateway

Johnjud-gateway is a routing and request handling service for the Johnjud project.

### What is Johnjud?
Johnjud is a pet adoption web application of the [CUVET For Animal Welfare Club](https://www.facebook.com/CUVETforAnimalWelfareClub)

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
2. Copy every `config.example.yaml` in `config` and paste it in the same directory with `.example` removed from its name.

3. Run `go mod download` to download all the dependencies.

### Running
1. Run `docker-compose up -d`
2. Run `make server` or `go run ./src/.`

### Testing
1. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`

## Other microservices/repositories of Johnjud
-  [Johnjud-gateway](https://github.com/isd-sgcu/johnjud-gateway): Routing and request handling
-  [Johnjud-auth](https://github.com/isd-sgcu/johnjud-auth): Authentication and authorization
-  [Johnjud-backend](https://github.com/isd-sgcu/johnjud-backend): Main business logic
-  [Johnjud-file](https://github.com/isd-sgcu/johnjud-file): File management service
- [Johnjud-proto](https://github.com/isd-sgcu/johnjud-proto): Protobuf files generator
- [Johnjud-go-proto](https://github.com/isd-sgcu/johnjud-go-proto): Generated protobuf files for golang
-  [Johnjud-frontend](https://github.com/isd-sgcu/johnjud-frontend): Frontend web application
