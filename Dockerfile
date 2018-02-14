FROM golang:1.9.4

ENV GOPATH $GOPATH:/go

RUN apt-get update

RUN mkdir -p /go/src/github.com/voidsatisfaction/haedal-project
ADD . /go/src/github.com/voidsatisfaction/haedal-project
WORKDIR /go/src/github.com/voidsatisfaction/haedal-project

CMD go run pilot/example/main.go
