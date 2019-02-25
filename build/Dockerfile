# 1. stage
FROM golang:alpine AS builder

# Install package manager
RUN apk --update add curl git

# Installing dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/wutchzone/auth/

COPY . .

# Install all packages
RUN go dep

RUN go run go build ./cmd/auth/*.go -o ./release/app

# 2. stage
FROM alpine

WORKDIR /root/

COPY --from=builder /go/src/github.com/wutchzone/auth/release .

CMD ["./app"]
