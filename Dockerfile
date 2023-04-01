FROM golang:1.22-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o


