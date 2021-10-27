.PHONY: build run

build: 
	go build -o bin/app main.go
	chmod 777 ./bin/app

run: build
	./bin/app