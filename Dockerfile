FROM golang:tip-alpine3.24 AS builder
RUN mkdir /app
WORKDIR /app

ADD . .
RUN go build

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder  /app/fake-names .
COPY --from=builder  /app/names.csv .

ENTRYPOINT ["./fake-names"]
