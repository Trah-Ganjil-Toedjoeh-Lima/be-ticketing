FROM golang:1.19-alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY=https://goproxy.io,direct

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base

RUN apk update
RUN apk add vips-dev

WORKDIR /go/src
COPY . .
RUN go mod download -x

RUN go build -o ./bin/app ./cmd/app/main.go
RUN go build -o ./bin/migrator ./cmd/migrator/main.go
RUN go build -o ./bin/email ./cmd/email/main.go

FROM alpine:latest AS runner

RUN apk update
RUN apk add vips-dev
RUN apk add terminus-font font-inconsolata font-dejavu font-noto font-noto-cjk font-awesome font-noto-extra

WORKDIR /ticketing-gmcgo

COPY --from=builder /go/src/bin /ticketing-gmcgo
COPY .env /ticketing-gmcgo

RUN mkdir "/ticketing-gmcgo/storage"
RUN mkdir "/ticketing-gmcgo/resource"


EXPOSE 8080 8080
CMD ["/ticketing-gmcgo/app"]
