FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/github.com/xtaci/trader
RUN go install github.com/xtaci/trader
ENTRYPOINT ["/go/bin/trader"]
EXPOSE 8888
