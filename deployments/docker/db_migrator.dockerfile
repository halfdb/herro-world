FROM golang:1.17-alpine3.14

RUN mkdir -p /go/src/github.com/halfdb/herro-world
WORKDIR /go/src/github.com/halfdb/herro-world
COPY . .

RUN apk add build-base
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.3.1

CMD make db-up
