.PHONY: build

LDFLAGS=-ldflags "-X=main.version=$(shell git describe --tags)"

release: build-release installer

install:
	go install $(LDFLAGS)

installer:
	godownloader --repo=unleashable/apker > ./install.sh

build-release: clean
	goreleaser release --rm-dist

build: clean
	go build $(LDFLAGS) -o bin/apker

clean:
	rm -rf bin/
	rm -rf dist/
