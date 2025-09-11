# ---------- Build Stage ----------
FROM golang:1.25-alpine3.22 AS builder

# 빌드 환경 준비
WORKDIR /app

# 의존성 캐시 최적화
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사
COPY . .

# 빌드
RUN go build -o main main.go

# ---------- Run Stage ----------
FROM alpine:3.22

# 실행 환경 준비
WORKDIR /app

# builder 스테이지에서 빌드된 실행파일만 복사
COPY --from=builder /app/main .

# 환경변수 (필요시)
ENV GIN_MODE=release

# 환경 변수 파일 복사
COPY app.env /app/app.env

# API 서버 포트
EXPOSE 8080

# 실행 시 환경 변수 로드
CMD ["sh", "-c", "source /app/app.env && /app/main"]
