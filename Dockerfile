FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd

COPY internal ./internal

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /norns

CMD ["/norns"]