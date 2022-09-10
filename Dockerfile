# build stage
FROM golang:1.19-alpine AS builder

LABEL maintainer="Joel Santos <joe@joesantos.io>"

WORKDIR /app
COPY . .

RUN apk update && apk add --virtual build-dependencies build-base gcc git make

RUN go mod download
RUN go build -v -o /app/book

# final stage
FROM alpine

ENV DB_PATH=/data/data.sqlite

WORKDIR /app/

RUN mkdir -p /data/user
COPY --from=builder /app/book /app/book

CMD ["/app/book", "/data/user/input.csv", "/data/user/output.csv"]