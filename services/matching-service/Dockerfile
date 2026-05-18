FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /app/matching-api \
  ./services/matching-service/cmd/matching-service

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/matching-api ./matching-api

RUN chmod +x ./matching-api

ENTRYPOINT ["./matching-api"]
