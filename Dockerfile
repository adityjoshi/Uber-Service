ARG SERVICE

# ─── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

ARG SERVICE
WORKDIR /app

COPY ${SERVICE}/go.mod ${SERVICE}/go.sum ./
RUN go mod download

COPY ${SERVICE}/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/${SERVICE}/

# ─── Final stage ─────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/service /bin/service

ENTRYPOINT ["/bin/service"]

