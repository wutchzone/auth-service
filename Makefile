.PHONY: build dev

default:
	@ echo "Available commands:"
	@ echo "make build 	- builds image"
	@ echo "make dev 	- runs dev docker-compose"

build:
	@ echo "Building image"
	docker image build --tag "wutchzone/auth" -f build/Dockerfile .

build:
    @ echo "Starting docker compose"
