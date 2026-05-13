FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /app/ride-api \
  ./services/ride-service/cmd/ride-api

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/ride-api ./ride-api

RUN chmod +x ./ride-api

ENTRYPOINT ["./ride-api"]
