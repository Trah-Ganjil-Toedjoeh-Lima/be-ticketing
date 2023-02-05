FROM golang:1.19-alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base

COPY go.mod go.sum /go/src/github.com/frchandra/ticketing-gmcgo/
WORKDIR /go/src/github.com/frchandra/ticketing-gmcgo
RUN apk update
RUN apk add vips-dev
RUN
RUN go mod download -x
#TODO modify qrcode library
COPY . /go/scr/github.com/frchandra/ticketing-gmcgo
RUN go build -o /usr/bin/ticketing-gmcgo/migrator github.com/frchandra/ticketing-gmcgo/cmd/migrator
RUN go build -o /usr/bin/ticketing-gmcgo github.com/frchandra/ticketing-gmcgo/cmd/app

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/ticketing-gmcgo"]

