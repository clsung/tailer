FROM golang:1.5.3-alpine

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

RUN apk add --update git
RUN go get -u github.com/clsung/tailer/cmd/tailer
RUN mkdir -p /tmp

CMD ${GOPATH}/bin/tailer /tmp
