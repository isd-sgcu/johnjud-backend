proto:
	go get github.com/isd-sgcu/johnjud-go-proto@latest

publish:
	cat ./token.txt | docker login --username isd-team-sgcu --password-stdin ghcr.io
	docker build . -t ghcr.io/isd-sgcu/johnjud-backend
	docker push ghcr.io/isd-sgcu/johnjud-backend

test:
	go vet ./...
	go test  -v -coverpkg ./src/app/... -coverprofile coverage.out -covermode count ./src/app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	. ./tools/export-env.sh ; go run ./src/.
