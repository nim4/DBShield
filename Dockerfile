FROM golang:1.6.3
MAINTAINER Nima Ghotbi <ghotbi.nima@gmail.com>

ENV GOPATH /go

COPY . /go/src/github.com/nim4/DBShield
WORKDIR /go/src/github.com/nim4/DBShield
COPY conf/dbshield.yml /etc/dbshield.yml

RUN go get
RUN go build
ENTRYPOINT /go/bin/DBShield
