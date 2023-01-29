FROM golang:1.20-rc-bullseye AS builder
ENV GO111MODULE=on
WORKDIR /go/src/github.com/rafavaliev/todo
COPY . .
RUN make
RUN ls -la /go/src/github.com/rafavaliev/todo

FROM ubuntu:bionic-20190424
RUN apt-get update && apt-get install -y apt-utils ca-certificates

COPY --from=builder /go/src/github.com/rafavaliev/todo/main /todo