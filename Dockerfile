# Go 1.25 (필요 버전에 맞게 수정)
FROM golang:1.25-alpine AS builder

WORKDIR /app

# go.mod / go.sum 먼저 복사해서 의존성 캐시
COPY go.mod ./
# go.sum 있으면 같이
# COPY go.sum ./
RUN go mod download

# 나머지 소스 복사
COPY . .

# 바이너리 빌드
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-fill-service ./cmd/server/main.go

# 최소 런타임 이미지
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/order-fill-service .

EXPOSE 8080
ENTRYPOINT ["./order-fill-service"]
