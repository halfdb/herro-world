FROM golang:1.17-alpine3.14 AS builder

RUN mkdir -p /go/src/github.com/halfdb/herro-world
WORKDIR /go/src/github.com/halfdb/herro-world
COPY . .

RUN apk add build-base curl
RUN make install-tools
RUN make server

FROM alpine:3.14 AS production

COPY --from=builder /go/src/github.com/halfdb/herro-world/bin/herro-world /herro-world

CMD /herro-world
