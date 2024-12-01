# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Go modulesの依存関係をキャッシュするためにgo.modとgo.sumをコピー
COPY go.mod .
COPY go.sum .
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# マイグレーションファイルをコピー
COPY migrations ./migrations

# バイナリをビルド
RUN go build -o weather-bot ./cmd/main.go

# 実行用の軽量イメージにコピー
FROM alpine:latest

WORKDIR /root/

# 必要なパッケージをインストール
RUN apk add --no-cache ca-certificates

# ビルドされたバイナリとマイグレーションファイルをコピー
COPY --from=builder /app/weather-bot .
COPY --from=builder /app/migrations ./migrations

# 実行
CMD ["./weather-bot"]
