proto:
	go get github.com/isd-sgcu/johnjud-go-proto@latest

publish:
	cat ./token.txt | docker login --username isd-team-sgcu --password-stdin ghcr.io
	docker build . -t ghcr.io/isd-sgcu/johnjud-gateway
	docker push ghcr.io/isd-sgcu/johnjud-gateway

test:
	go vet ./...
	go test  -v -coverpkg ./src/app/... -coverprofile coverage.out -covermode count ./src/app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./src/.

mockGen:
	mockgen -source ./src/pkg/service/auth/auth.service.go -destination ./src/mocks/service/auth/auth.mock.go
	mockgen -source ./src/pkg/service/user/user.service.go -destination ./src/mocks/service/user/user.mock.go
	mockgen -source ./src/app/validator/validator.go -destination ./src/mocks/validator/validator.mock.go
	mockgen -source ./src/app/router/context.go -destination ./src/mocks/router/context.mock.go
