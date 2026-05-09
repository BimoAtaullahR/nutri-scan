FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/nutriscan-api ./cmd/api

FROM alpine:3.21

RUN addgroup -S nutriscan && adduser -S nutriscan -G nutriscan

USER nutriscan
WORKDIR /app

COPY --from=build /out/nutriscan-api /app/nutriscan-api

EXPOSE 8080

ENTRYPOINT ["/app/nutriscan-api"]
