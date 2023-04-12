# Step 1
FROM golang:alpine AS builder

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

# Step 2
FROM alpine:latest AS runner
RUN apk update
RUN apk add vips-dev
RUN apk add terminus-font font-inconsolata font-dejavu font-noto font-noto-cjk font-awesome font-noto-extra

WORKDIR /ticketing-gmcgo

COPY --from=builder /go/src/bin /ticketing-gmcgo
COPY storage /ticketing-gmcgo/storage
COPY resource /ticketing-gmcgo/resource

EXPOSE 5000
CMD ["/ticketing-gmcgo/app"]
