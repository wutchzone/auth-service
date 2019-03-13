.PHONY: build

default:
	@ echo "Available commands:"
	@ echo "make build 	- builds image"

build:
	@ echo "Building image"
	docker image build --tag "wutchzone/auth-service" -f build/Dockerfile .
