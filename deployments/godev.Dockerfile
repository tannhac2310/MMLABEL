FROM golang:1.21.2 AS dev_env

RUN apt-get update && apt-get install -y --no-install-recommends curl make git unzip

RUN go install github.com/go-delve/delve/cmd/dlv@v1.9.1

RUN curl -OL https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz
RUN tar zxvf modd-0.8-linux64.tgz -C /
RUN mv /modd-0.8-linux64/modd /usr/local/bin/modd

ENV GOPROXY="https://goproxy.io,direct"

COPY ./go.mod /
COPY ./go.sum /
RUN go mod download
