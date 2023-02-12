
build:
	go build -v ./...

test:
	go test ./...

vet: 
	go vet ./...

run:
	go run .

upgrade:
	go get -u