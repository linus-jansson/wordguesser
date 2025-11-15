# ---------- Build stage ----------
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY wordguesser-go/go.mod ./

RUN go mod download
COPY wordguesser-go/*.go ./

COPY words/valid-words.json ./valid-words.json

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wordguesser .

# ---------- Runtime stage ----------
FROM alpine:latest

WORKDIR /app

# Copy the binary and the word list from the builder
COPY --from=builder /app/wordguesser .
COPY --from=builder /app/valid-words.json .

# Env var default (overridden by docker-compose/.env)
ENV API_URL="https://ordel.se/play"

ENTRYPOINT ["./wordguesser"]
