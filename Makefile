.PHONY: parse build run clean

VERSION ?= dev

parse:
	go run ./cmd/parser

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/ghofig ./cmd/ghofig

run: build
	./bin/ghofig

clean:
	rm -rf bin/
