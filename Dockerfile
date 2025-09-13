# ---------- Build Stage ----------
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# ---------- Run Stage ----------
FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migration ./migration
COPY start.sh .

RUN chmod +x /app/start.sh

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
