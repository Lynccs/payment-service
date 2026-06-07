FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/payment-service


FROM alpine

WORKDIR /app

# Install golang-migrate
RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/ \
    && chmod +x /usr/local/bin/migrate

COPY --from=builder /app/payment-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/web ./web

EXPOSE 8082

CMD ["./payment-service"]