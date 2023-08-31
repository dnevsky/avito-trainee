FROM golang:1.20-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o avito-trainee cmd/main.go
 
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/avito-trainee /app/avito-trainee
COPY --from=builder /app/configs /app/configs

WORKDIR /app

EXPOSE 8002


ENV TZ Europe/Moscow
CMD ["./avito-trainee"]