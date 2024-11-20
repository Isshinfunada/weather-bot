# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Go modulesの依存関係をキャッシュするためにgo.modとgo.sumをコピー
COPY go.mod .
COPY go.sum .
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# バイナリをビルド
RUN go build -o weather-bot ./cmd/main.go

# 実行用の軽量イメージにコピー
FROM alpine:latest

WORKDIR /root/

# ビルドされたバイナリをコピー
COPY --from=builder /app/weather-bot .

# 実行
CMD ["./weather-bot"]
