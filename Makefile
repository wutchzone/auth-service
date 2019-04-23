.PHONY: all build test buildimage

all:
	@ echo "Available commands:"
	@ echo "make build 	- builds image"

build:
    @ Building application
    @ go build -o ./release/app ./cmd/auth/*.go

test:
    @ Run tests

buildimage:
	@ echo "Building image"
	docker image buildimage --tag "wutchzone/auth-service" -f build/Dockerfile .
