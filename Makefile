.PHONY: build

LDFLAGS=-ldflags "-X=main.version=$(shell git describe --tags)"

build: clean
	go build $(LDFLAGS) -o bin/apker

install:
	go install $(LDFLAGS)

release: clean
	env GOOS=linux go build $(LDFLAGS) -o bin/apker
	env GOOS=darwin go build $(LDFLAGS) -o bin/apker-darwin
	env GOOS=windows go build $(LDFLAGS) -o bin/apker.exe

clean:
	rm -rf bin/
