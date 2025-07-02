FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/service ./cmd/

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/service /usr/local/bin/service
COPY --from=builder /app/web/static /usr/local/share/static

EXPOSE 8081

CMD ["service"]
