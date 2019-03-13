# 1. stage
FROM golang:alpine AS builder

# Install package manager
RUN apk --update add curl git

# Installing dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/wutchzone/auth-service/

COPY . .

# Install all packages
RUN dep ensure

RUN go build -o ./release/app ./cmd/auth/*.go

# 2. stage
FROM alpine

WORKDIR /root/

COPY --from=builder /go/src/github.com/wutchzone/auth-service/release .

COPY ./configs/config.json .

CMD ["./app", "./config.json"]
