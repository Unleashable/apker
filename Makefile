.PHONY: build

LDFLAGS=-ldflags "-X=main.version=$(shell git describe --tags)"

build: clean
	go build $(LDFLAGS) -o bin/apker

install:
	go install $(LDFLAGS)

clean:
	rm -rf bin/

