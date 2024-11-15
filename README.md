## プロジェクト概要

このプロジェクトは、ユーザーが指定した時刻と地域に基づいて天気情報をLINEを通じて通知するLINE BOTの開発を目的としている。LINEアカウントを利用した認証機能を活用し、シンプルで効率的な天気情報通知システムを構築する。

## 機能概要

- **天気通知機能**: ユーザーが指定した時刻および地域に基づき、天気情報をLINEで通知する。
- **ユーザー認証**: LINEアカウント情報を利用したユーザー認証。

## 技術スタック

### 1. アプリケーションアーキテクチャ

- **コンテナ技術**: Kubernetesを利用してコンテナを管理。HelmやKustomizeを使用してKubernetesのマニフェストを効率的に管理し、後々の運用も見据えた設定を行う。
- **API連携**: 以下のAPIを活用して天気情報と位置情報を取得する。
  - 天気予報: OpenWeatherMap API
  - 緯度経度取得用: Google Geocoding API

### 2. システムインフラ

- **CI/CD**:
  - **CI**: Github Actionsにより自動テストとデプロイを効率化。
  - **CD**: ArgoCDを採用し、Kubernetes上での継続的デリバリーを実現。
- **デプロイ**: Fly.ioを利用してアプリケーションをホスティング。
- **Infrastructure as Code**: Terraformを利用してインフラをコードで管理。

### 3. 認証とメッセージ送信

- **ユーザー認証**: LINEアカウントを利用してユーザー認証を行う。
- **メッセージ通知**: LINE Messaging APIを使って天気情報を通知。

### 4. データ管理

- **データベース**:
  - **キャッシュ**: Redisを利用。
  - **RDBMS**: PostgreSQLを使用。
- **SQL生成ツール**: sqlcを利用して型安全なSQLクエリを作成。

### 5. セキュリティ

- **APIキー管理**: APIキーは環境変数として管理し、セキュリティを確保。
- **データ保護**: ユーザー情報の暗号化と、ネットワーク通信のHTTPS/TLSによる暗号化を徹底。

## 今後のステップ

1. **Kubernetes導入計画**: 現在はDockerを使用しているが、Kubernetesの導入に向けて負荷分散やCI/CDの統合計画を具体化。
2. **スケーラビリティテスト**: 負荷をシミュレートしてシステムのボトルネックを特定し、スケーラビリティを確認。

## 開発環境のセットアップ

1. リポジトリをクローン
   ```bash
   git clone <repository-url>
   ```
2. Dockerコンテナをビルドし起動
   ```bash
   docker-compose up --build
   ```
3. 必要な環境変数を設定
   - `OPENWEATHERMAP_API_KEY`: OpenWeatherMap APIのキー。
   - `GOOGLE_GEOCODING_API_KEY`: Google Geocoding APIのキー。

## ライセンス

このプロジェクトはMITライセンスの下で公開されている。

# ディレクトリ構成
クリーンアーキテクチャを目指したい。

/project-root
  ├── /cmd                 # アプリケーションのエントリーポイント（main.goなど）
  ├── /internal            # アプリケーション固有の内部ロジック
  │     ├── /entity        # ドメインモデル（エンティティ）
  │     │     └── user.go  # ユーザーエンティティ（例：LINEユーザー情報など）
  │     ├── /usecase       # ユースケース
  │     │     ├── notify_weather.go  # 天気通知のユースケース
  │     │     └── auth_user.go       # LINE認証処理のユースケース
  │     ├── /interfaces    # インターフェースアダプタ（データ変換やリポジトリの実装）
  │     │     ├── /controller # コントローラ（HTTPリクエストを処理する部分）
  │     │     │     └── line_webhook.go # LINEからのWebhook処理
  │     │     ├── /repository # データベースリポジトリ
  │     │     │     └── user_repository.go # ユーザーデータの操作
  │     │     └── /presenter # 出力変換、APIレスポンスの整形など
  │     │           └── line_presenter.go # LINE向けのメッセージ整形
  ├── /pkg                 # 再利用可能なコード（汎用的なライブラリ）
  │     ├── /weather       # 天気情報取得のロジック
  │     ├── /geo           # 緯度経度取得のロジック
  │     └── /lineapi       # LINE Messaging APIのクライアントラッパー
  ├── /configs             # 設定ファイル（YAML, JSONなど）
  ├── /scripts             # デプロイやビルドなどのスクリプト
  ├── /build               # コンテナやCI/CDで使用するビルド関連のファイル
  └── /test                # テストコード
        ├── /unit          # ユニットテスト
        └── /integration   # 統合テスト