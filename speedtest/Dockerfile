# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd

EXPOSE 8081

CMD [ "/app/server" ] 