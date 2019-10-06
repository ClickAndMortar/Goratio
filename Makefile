all: build

default: build

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/goratio-linux main.go
	GOOS=darwin GOARCH=amd64 go build -o ./bin/goratio-darwin main.go

docker:
	docker build -t clickandmortar/goratio .

test:
	go test -coverprofile=coverage.out
