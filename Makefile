all: package

test:
	docker run --rm -v "${PWD}/src":/usr/src/weightServer -w /usr/src/weightServer golang:1.6 go test -v

build:
	cd src && go build -o weightServer

package: build
	docker build -t nwik/weightserver .

