FROM golang:1.19-bullseye

COPY go.mod go.sum /go/src/github.com/frchandra/ticketing-gmcgo/
WORKDIR /go/src/github.com/frchandra/ticketing-gmcgo
RUN apt-get update
RUN apt-get install -y libvips-dev
RUN go mod download -x
#TODO modify qrcode library
COPY . /go/scr/github.com/frchandra/ticketing-gmcgo
RUN go build -o /usr/bin/ticketing-gmcgo/migrator github.com/frchandra/ticketing-gmcgo/cmd/migrator
RUN go build -o /usr/bin/ticketing-gmcgo github.com/frchandra/ticketing-gmcgo/cmd/app

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/ticketing-gmcgo"]

