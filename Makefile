.PHONY: build

LDFLAGS=-ldflags "-X=main.version=$(shell git describe --tags)"

build: clean
	go build $(LDFLAGS) -o bin/apker

release: buildrelease installer

install:
	chmod +x bin/apker
	mv bin/apker /usr/bin/apker
	install -C autocomplete/bash /usr/share/bash-completion/completions/apker

installer:
	godownloader --repo=unleashable/apker > ./install.sh

buildrelease: clean
	goreleaser release --rm-dist

clean:
	rm -rf bin/
	rm -rf dist/
