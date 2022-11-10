# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine
WORKDIR /conBot

COPY go.mod ./
COPY go.sum ./
COPY startContainer.sh /

RUN go mod download
COPY *.go ./

RUN go build -o /bot

ENV DATADIR=/con

CMD [ "/startContainer.sh" ]