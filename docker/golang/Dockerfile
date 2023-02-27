FROM golang:1.19.6 AS builder

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app
RUN go mod tidy & go build -o /main_server

FROM golang:1.19.6-buster

COPY --from=builder /go/src/app/.env /.env
COPY --from=builder /main_server /main_server
