FROM golang:latest

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app
RUN go mod tidy & go build -o /main_server
