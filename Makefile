.PHONY : format install build

run:
	go run ./bin/app/main.go

format:
	gofmt -s -w .

run-this:
	echo "hello"

everything-oke:
	go run ./bin/app/main.go

install:
	go mod download

build:
	go build -tags musl -o main ./bin/app

start:
	./main

# live reload using nodemon: npm -g i nodemon
run-nodemon:
	nodemon --exec go run ./bin/app/main.go
