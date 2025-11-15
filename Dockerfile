FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache ca-certificates gzip

WORKDIR /app

COPY wordguesser-go/go.mod ./

COPY wordguesser-go/*.go ./
COPY words/valid-words.json ./valid-words.json
# We compress the word list and embed it in build time
RUN gzip -c valid-words.json > valid-words.json.gz
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o wordguesser .
RUN rm ./valid-words.json.gz valid-words.json

# Runtime stage - using scratch for minimal image size
FROM scratch

WORKDIR /app

COPY --from=builder /app/wordguesser .
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/app/wordguesser"]
