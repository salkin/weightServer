all: package

build:
	go build -o weightServer

package: build
	docker build -t nwik/weightserver .

