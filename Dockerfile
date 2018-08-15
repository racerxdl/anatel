FROM golang:alpine

RUN apk update

RUN apk add git ca-certificates

ADD . /go/src/github.com/racerxdl/anatel

WORKDIR /go/src/github.com/racerxdl/anatel

RUN go get -v && go build -o anatel_worker

ENTRYPOINT /go/src/github.com/racerxdl/anatel/anatel_worker