FROM golang:latest

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
RUN go mod tidy \
  && go install github.com/cweill/gotests/gotests@latest \
  && go install github.com/fatih/gomodifytags@latest \
  && go install github.com/josharian/impl@latest \
  && go install github.com/haya14busa/goplay/cmd/goplay@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install honnef.co/go/tools/cmd/staticcheck@latest \ 
  && go install golang.org/x/tools/gopls@latest
