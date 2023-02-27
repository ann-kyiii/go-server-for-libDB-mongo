FROM golang:1.19.6 AS builder

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY ./add_borrower /go/src/app
COPY ./.env /go/src/app/.env
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
RUN go mod tidy & go build -o /add_borrower

FROM mongo:6.0.4

COPY --from=builder /add_borrower /add_borrower
