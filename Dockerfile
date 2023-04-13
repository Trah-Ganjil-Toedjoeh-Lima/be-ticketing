# Step 1
FROM golang:latest AS builder
ENV GOPROXY=https://goproxy.io,direct

RUN apt update && apt install -y apt-utils
RUN apt install --no-install-recommends -y libvips-dev

WORKDIR /go/src
COPY . .
RUN go mod download -x

RUN CGO_ENABLED=1 go build -o ./bin/app ./cmd/app/main.go
RUN CGO_ENABLED=1 go build -o ./bin/migrator ./cmd/migrator/main.go
RUN CGO_ENABLED=1 go build -o ./bin/email ./cmd/email/main.go

# Step 2
FROM debian:stable-slim AS runner
RUN apt update && apt install -y apt-utils
RUN apt install -y --no-install-recommends libvips xfonts-terminus fonts-inconsolata fonts-dejavu fonts-noto fonts-noto-cjk fonts-font-awesome fonts-noto-extra openssl ca-certificates

WORKDIR /ticketing-gmcgo

COPY --from=builder /go/src/bin /ticketing-gmcgo
COPY storage /ticketing-gmcgo/storage
COPY resource /ticketing-gmcgo/resource

EXPOSE 5000
CMD ["/ticketing-gmcgo/app"]
