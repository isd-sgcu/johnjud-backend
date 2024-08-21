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
	. ./tools/export-env.sh ; go run ./cmd/.

mock-gen:
	mockgen -source ./internal/auth/auth.service.go -destination ./mocks/service/auth/auth.mock.go
	mockgen -source ./internal/pkg/service/user/user.service.go -destination ./mocks/service/user/user.mock.go
	mockgen -source ./internal/pkg/service/pet/pet.service.go -destination ./mocks/service/pet/pet.mock.go
	mockgen -source ./internal/pkg/service/image/image.service.go -destination ./mocks/service/image/image.mock.go
	mockgen -source ./internal/validator/validator.go -destination ./mocks/validator/validator.mock.go
	mockgen -source ./internal/router/context.go -destination ./mocks/router/context.mock.go

create-doc:
	swag init -d ./src -o ./src/docs -md ./src/docs/markdown