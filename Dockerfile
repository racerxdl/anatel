FROM golang:alpine as build

RUN apk update

RUN apk add git ca-certificates

ADD . /go/src/github.com/racerxdl/anatel

WORKDIR /go/src/github.com/racerxdl/anatel

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o anatel_worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/racerxdl/anatel/anatel_worker .

CMD ["./anatel_worker"]