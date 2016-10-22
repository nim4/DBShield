FROM golang:1.6.3
MAINTAINER Nima Ghotbi <ghotbi.nima@gmail.com>

ENV GOPATH /go

COPY . /go/src/github.com/nim4/DBShield
WORKDIR /go/src/github.com/nim4/DBShield
COPY conf/dbshield.yml /etc/dbshield.yml

RUN openssl genrsa -out cert/server-key.pem 2048
RUN openssl req -new -x509 -sha256 -key cert/server-key.pem -out cert/server-cert.pem -days 3650 -subj '/CN=DBShield/O=DBShield/C=TR'

RUN go get
RUN go build
ENTRYPOINT /go/bin/DBShield
