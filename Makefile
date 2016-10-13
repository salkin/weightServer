all: package

build:
	cd src && go build -o weightServer

package: build
	docker build -t nwik/weightserver .

