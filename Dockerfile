# ---------- Build Stage ----------
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go


# ---------- Run Stage ----------
FROM alpine:3.22
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration
COPY start.sh .
COPY wait-for.sh .

RUN chmod +x /app/start.sh

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
