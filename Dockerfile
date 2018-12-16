FROM golang:1.11.3-stretch

RUN go get "golang.org/x/tools/cmd/goimports" && go get "github.com/gopherjs/gopherjs"