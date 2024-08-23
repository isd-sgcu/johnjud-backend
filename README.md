# Johnjud-backend

Johnjud-backend is a routing and request handling service for the Johnjud project.

### What is Johnjud?
Johnjud is a pet adoption web application of the [CUVET For Animal Welfare Club](https://www.facebook.com/CUVETforAnimalWelfareClub)

## Stack

-   golang
-   go-fiber

## Getting Started

### Prerequisites

-   golang 1.21 or [later](https://go.dev)
-   docker
-   makefile

### Installation

1. Clone this repo
2. Copy `.env.template` in root directory and paste it in the same directory as `.env` with proper values.
3. Run `go mod download` to download all the dependencies.

### Running
1. Run `docker-compose -f docker-compose.example.yaml up`
2. Run `make server` or `go run ./cmd/.`

### Testing
1. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`

## Other repositories of Johnjud
- [Johnjud-backend](https://github.com/isd-sgcu/johnjud-backend)
- [Johnjud-frontend](https://github.com/isd-sgcu/johnjud-frontend)
