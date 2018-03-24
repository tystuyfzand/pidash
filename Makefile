install-pre:
	go get github.com/Masterminds/glide
	gem install --no-ri --no-rdoc fpm

install-deps:
	glide install

build:
	go build -o pidash pidash.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./pidash ${DESTDIR}/usr/local/bin/glide

test: install-pre install-deps
	go test ./dashboard/module

build-armv7:
	GOOS=linux GOARCH=arm go build -o dist/armv7/pidash pidash.go
	ARCH=armv7 bash packaging/build-package.sh

build-amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/amd64/pidash pidash.go
	ARCH=amd64 bash packaging/build-package.sh

build-all: build-armv7 build-amd64

.PHONY: install-pre install-deps test