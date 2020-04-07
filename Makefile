.PHONY: build

LDFLAGS=-ldflags "-X=main.version=$(shell git describe --tags)"

build: clean
	go build $(LDFLAGS) -o bin/apker

install:
	go install $(LDFLAGS)

release: clean
	goreleaser release --rm-dist

installer:
	godownloader --repo=unleashable/apker > ./install.sh

clean:
	rm -rf bin/
