.PHONY: build run

build:
	mkdir -p bin
	go build -o bin/dingo

run:
	go run . --logmode human
