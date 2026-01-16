.PHONY: parse build run clean

parse:
	go run ./cmd/parser

build:
	go build -o bin/ghofig ./cmd/ghofig

run: build
	./bin/ghofig

clean:
	rm -rf bin/ data/ghofig.db
