proto:
	go get github.com/isd-sgcu/johnjud-go-proto@latest

publish:
	cat ./token.txt | docker login --username isd-team-sgcu --password-stdin ghcr.io
	docker build . -t ghcr.io/isd-sgcu/johnjud-gateway
	docker push ghcr.io/isd-sgcu/johnjud-gateway

test:
	go vet ./...
	go test  -v -coverpkg ./internal/... -coverprofile coverage.out -covermode count ./internal/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./cmd/.

docker-qa:
	docker-compose -f docker-compose.qa.yaml up

mock-gen:
	mockgen -source ./internal/cache/cache.repository.go -destination ./mocks/repository/cache/cache.mock.go
	mockgen -source ./internal/auth/auth.repository.go -destination ./mocks/repository/auth/auth.mock.go
	mockgen -source ./internal/auth/auth.service.go -destination ./mocks/service/auth/auth.mock.go
	mockgen -source ./internal/user/user.service.go -destination ./mocks/service/user/user.mock.go
	mockgen -source ./internal/pet/pet.service.go -destination ./mocks/service/pet/pet.mock.go
	mockgen -source ./internal/image/image.service.go -destination ./mocks/service/image/image.mock.go
	mockgen -source ./internal/validator/validator.go -destination ./mocks/validator/validator.mock.go
	mockgen -source ./internal/router/context.go -destination ./mocks/router/context.mock.go

create-doc:
	swag init -d ./internal -g ../cmd/main.go -o ./docs -md ./docs/markdown --parseDependency --parseInternal