FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd

COPY internal ./internal

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /norns

# Deployment to lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /norns /norns

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/norns"]