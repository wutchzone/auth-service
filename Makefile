.PHONY: all build test buildimage

all:
	@ echo "Available commands:"
	@ echo "make build 	- builds image"

build:
	@ Building application
	@ go build -o ./release/app ./cmd/auth/*.go

test:
	@ echo Run tests
	go test ./pkg/accountdb
	go test ./pkg/sessiondb

dev:
	docker-compose -f ./deployments/docker-compose.yml up -d

clean:
	docker-compose -f ./deployments/docker-compose.yml down

buildimage:
	@ echo "Building image"
	docker image buildimage --tag "wutchzone/auth-service" -f build/Dockerfile .
