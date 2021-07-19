FROM golang:rc-alpine AS builder
RUN mkdir /app
WORKDIR /app

ADD . .
RUN go build

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder  /app/fake-names .

ENTRYPOINT ["./fake-names"]
