# ビルドステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Goモジュールをコピーして依存関係をインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN go build -o weather-bot ./cmd/main.go

# 実行ステージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk add --no-cache ca-certificates

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドされたバイナリをコピー
COPY --from=builder /app/weather-bot .

# マイグレーションスクリプトをコピー
COPY --from=builder /app/db ./db

# アプリケーションのポートを公開
EXPOSE 8080

# 実行
CMD ["./weather-bot"]
