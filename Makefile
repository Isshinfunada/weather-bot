.PHONY: build docker-up docker-down migrate-local migrate-k8s deploy start

# Dockerイメージのビルド
build:
	docker build -t isshinfunada/weather-bot:latest .

# ローカルでDocker Composeを使ってサービスを起動
docker-up:
	docker compose up -d

# ローカルでDocker Composeを使ってサービスを停止
docker-down:
	docker compose down -v

# Kubernetes上でアプリケーションをデプロイ
deploy: build
	docker push isshinfunada/weather-bot:latest 
	helm upgrade --install weather-bot ./charts/weather-bot -n weather-bot
