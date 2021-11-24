FROM golang:1.17-alpine3.14

RUN mkdir -p /go/src/github.com/halfdb/herro-world
WORKDIR /go/src/github.com/halfdb/herro-world
COPY . .

RUN apk add build-base
RUN make server

CMD ./bin/herro-world
