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

RUN go build -o ./bin/app ./cmd/app

FROM alpine:latest AS runner
COPY --from=builder /go/src/bin/app /
CMD ["./bin/app"]
