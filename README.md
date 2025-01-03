## プロジェクト概要

このプロジェクトは、ユーザーが指定した主要都市に基づいて天気情報をLINEを通じて通知するLINE BOTの開発を目的としています。気象庁（JMA）が提供する天気予報データを活用し、シンプルで効率的な天気情報通知システムを構築します。

## 機能概要

- 天気通知機能: ユーザーが選択した主要都市の天気情報をLINEで通知する。
- ユーザー認証: LINEアカウント情報を利用したユーザー認証。
- 通知スケジュール管理: ユーザーごとに通知の時刻と都市を設定し、指定した時間に天気情報を送信。 
- 都市選択機能: ユーザーが初回設定時に主要都市を選択し、後から変更可能。

## 技術スタック

1. アプリケーションアーキテクチャ

- コンテナ技術: Kubernetesを利用してコンテナを管理。HelmやKustomizeを使用してKubernetesのマニフェストを効率的に管理し、後々の運用も見据えた設定を行う。
- API連携: 以下のAPIを活用して天気情報を取得する。
- 天気予報: 気象庁（JMA）API (https://www.jma.go.jp/bosai/forecast/data/forecast/{location_id}.json)

2. システムインフラ

- CI/CD:
    - CI: GitHub Actionsにより自動テストとデプロイを効率化。
    - CD: ArgoCDを採用し、Kubernetes上での継続的デリバリーを実現。
- デプロイ: 
    - Fly.ioを利用してアプリケーションをホスティング。
- IaC: 
    - Terraformを利用してインフラをコードで管理。

3. 認証とメッセージ送信

- ユーザー認証: 
    - LINEアカウントを利用してユーザー認証を行う。
- メッセージ通知: 
    - LINE Messaging APIを使って天気情報を通知。

4. データ管理

- データベース:
    - キャッシュ: Redisを利用。
    - RDBMS: PostgreSQLを使用。
    - SQL生成ツール: sqlcを利用して型安全なSQLクエリを作成。

5. セキュリティ

- APIキー管理: 
    - APIキーは環境変数として管理し、セキュリティを確保。
- データ保護: 
    - ユーザー情報の暗号化と、ネットワーク通信のHTTPS/TLSによる暗号化を徹底。

## 詳細設計

1. 天気通知機能

機能概要

ユーザーが選択した主要都市に基づき、天気情報をLINEで通知する機能です。基本的には雨の時のみ通知します。

実装のポイント

データ取得

- 気象庁APIの活用:
- 主要都市のlocation_idを事前に確認し、定期的に天気情報を取得します。
- APIのレスポンスデータを解析し、必要な情報（天気、気温、湿度、風速、降水確率など）を抽出します。
- APIの効率的な利用:
- CronJobを使用して1時間おきに全主要都市の天気情報を取得し、キャッシュ（例：Redis）に保存します。
- キャッシュを利用することで、ユーザーリクエスト時にリアルタイムでAPIを呼び出す必要を減らし、システムの効率を向上させます。

メッセージフォーマット

- LINE Messaging APIの利用:
- リッチメニューやテンプレートメッセージを活用して、視覚的にわかりやすい天気情報を提供します。
- メッセージの内容やレイアウトをユーザーの設定に応じてカスタマイズ可能にします。

レート制限対策

- APIのレート制限:
- 気象庁APIの利用頻度をCronJobで制御し、定期的にデータを取得することでレート制限を管理します。

スケーリングとパフォーマンス

- 非同期処理の導入:
- 通知の送信を非同期に処理することで、システムの応答性を向上させます。

実装における課題

- APIのレート制限:
- 気象庁APIの利用制限に注意し、必要に応じて取得頻度を調整。
- データの正確性と更新頻度:
- 天気情報の更新頻度に応じて、適切なデータ取得タイミングを設定します。

2. ユーザー認証

機能概要

LINEアカウント情報を利用してユーザー認証を行い、個別の設定や通知を管理します。

実装のポイント

認証フロー

- LINE Loginの利用:
- LINEのOAuth 2.0を使用してユーザー認証を実装します。
- ユーザーがLINEアカウントを介してログインし、必要な権限を取得します。

アクセストークンの管理

- アクセストークンの安全な保存:
- LINEから取得したアクセストークンを安全に保存・管理します。
- トークンの有効期限を監視し、必要に応じてリフレッシュトークンを使用して更新します。

ユーザーデータの管理

- データベース設計:
- ユーザー情報（LINEのユーザーID、アクセストークン、設定情報など）をPostgreSQLに保存します。
- セキュリティを考慮し、個人情報は暗号化して保存します。
- プライバシーとセキュリティ:
- データベースへのアクセス制御や暗号化を徹底します。

実装における課題

- 認証の信頼性:
- 認証フロー中のエラーハンドリングを適切に行い、ユーザーにわかりやすいエラーメッセージを提供します。
- スケーラビリティ:
- 多数のユーザーが同時に認証を行った場合でも、システムが安定して動作するように設計します。

3. 通知スケジュール管理

機能概要

ユーザーごとに通知の時刻と都市を設定し、指定した時間に天気情報を送信します。

実装のポイント

スケジューリングシステム

- ジョブスケジューラーの選定:
- CronJob: 主要都市の天気情報を定期的に取得し、キャッシュを更新します。
- 通常のJob: ユーザーごとの通知ジョブを管理し、指定された時間に天気情報を送信します。
- スケジュールの動的管理:
- ユーザーが通知時間や都市を変更した際に、スケジュールを動的に更新できる仕組みを構築します。

データ管理

- 通知設定の保存:
- ユーザーごとの通知時刻、都市、通知形式などの設定をデータベースに保存します。
- スケジュールジョブをトリガーするための情報も含めます。

タイムゾーンの考慮
- 基本的に日本国内を対象とします。
- すべてのユーザーのタイムゾーンを日本標準時（JST）に統一します。

実装における課題

- 時間の正確性:
	- サーバーのタイムゾーンとユーザーのタイムゾーンの違いを適切に管理し、通知が正確な時間に送信されるようにします。
- スケーラビリティとパフォーマンス:
	- 大規模なユーザー数に対応するため、スケジューリングシステムのパフォーマンスを最適化します。

4. 主要都市の天気情報管理

機能概要

主要都市の天気情報を事前に取得し、ユーザーが選択できるようにします。これにより、ユーザーが個別の位置情報を入力する必要がなくなり、プライバシーリスクを軽減します。

実装のポイント

主要都市の選定と管理

- 主要都市リストの作成:
- 全国の主要都市のlocation_idをリストアップします。例えば、東京、大阪、名古屋、札幌など。enumで管理します。
- データの取得と更新:
- CronJobを使用して、1時間おきに主要都市の天気情報を気象庁APIから取得し、キャッシュに保存します。

データ取得の自動化

- CronJobの設定:
- Kubernetes CronJobを設定し、定期的に天気情報を取得するジョブを実行します。
- キャッシュの活用:
- 取得した天気情報をRedisにキャッシュし、ユーザーリクエスト時に迅速に提供できるようにします。

ユーザーインターフェースの調整

- 都市選択のUI:
- 初回設定時にユーザーに主要都市を選択させるインターフェースを提供します。
- LINEのリッチメニューやテンプレートメッセージを活用し、ユーザーが簡単に都市を選択できるようにします。

実装における課題

- データの一貫性:
- 天気情報の取得タイミングとユーザーへの通知タイミングの整合性を保つ必要があります。
- エラーハンドリング:
- 天気情報の取得に失敗した場合の対応策（リトライ、通知しないなど）を実装します。

5. キャッシュ戦略

キャッシュのタイミング

- CronJobによる定期取得:
- Kubernetes CronJobを使用して、主要都市の天気情報を1時間おきに取得し、Redisにキャッシュします。

キャッシュの有効期限

- 天気情報の有効期限を1時間に設定し、CronJobで定期的に更新します。これにより、データの鮮度を維持しつつ、キャッシュの有効期限を適切に管理します。


6. データベース設計

テーブル構造

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    line_user_id VARCHAR(255) UNIQUE NOT NULL,
    selected_city_id VARCHAR(10) REFERENCES area_class20(id),
    notify_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

データベース設計のポイント

- ユーザー情報: users テーブルにLINEユーザーID、アクセストークン、選択した都市のlocation_idを保存。
- 都市情報: cities テーブルに主要都市のlocation_idと名前、地域を保存。
- 通知設定: notification_settings テーブルにユーザーごとの通知時刻、形式、タイムゾーンを保存。
- 通知履歴: notification_history テーブルに過去の通知データを保存し、履歴管理を実現。

7. 主要都市の天気情報取得ジョブ

CronJobの設定

```yaml
# build/kubernetes/jma_weather_cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: jm-weathe-cronjob
spec:
  schedule: "0 * * * *" # 毎時0分に実行
  jobTemplate:
    spec:
      template:
        spec:
          containers:
        - name: jm-weathe-fetcher
            image: you-docke-image
            env:
          - name: JMA_API_URL
              value: "https://www.jma.go.jp/bosai/forecast/data/forecast/"
          - name: CITIES_FILE
              value: "/configs/cities.json"
          - name: REDIS_ADDR
              value: "redis:6379"
          - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: redi-secret
                  key: password
          restartPolicy: OnFailure
```

ジョブの実装例

```go
// internal/usecase/fetch_weather.go
package usecase

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    ...
)

func (uc *FetchWeatherUseCase) Execute() error {
    // 主要都市リストの読み込み
    cities, err := uc.loadCities()
    if err != nil {
        return err
    }

    for _, city := range cities {
        weatherData, err := uc.fetchWeather(city.ID)
        if err != nil {
            fmt.Printf("Failed to fetch weather for %s: %v\n", city.Name, err)
            continue
        }

        // キャッシュに保存
        cacheKey := fmt.Sprintf("weather:%s", city.ID)
        weatherJSON, err := json.Marshal(weatherData)
        if err != nil {
            fmt.Printf("Failed to marshal weather data for %s: %v\n", city.Name, err)
            continue
        }

        err = uc.Cache.Set(context.Background(), cacheKey, string(weatherJSON), 1*time.Hour)
        if err != nil {
            fmt.Printf("Failed to set cache for %s: %v\n", city.Name, err)
        }
    }

    return nil
}
```

主要都市リストの管理

cities.jsonの例

```json
[
    {
        "id": "130000", // 東京
        "name": "東京",
        "region": "関東地方"
    },
    {
        "id": "270000", // 大阪
        "name": "大阪",
        "region": "近畿地方"
    },
    {
        "id": "140000", // 札幌
        "name": "札幌",
        "region": "北海道地方"
    }
    // その他の主要都市
]
```

ユーザーインターフェースの調整

都市選択のフロー

1. 初回設定時:
   - ユーザーがLINEを通じてボットと対話を開始。
   - ボットが主要都市のリストをリッチメニューやテンプレートメッセージで提示。
   - ユーザーが都市を選択し、選択した都市IDをデータベースに保存。
2. 後からの変更:
   - ユーザーがメニューから「都市を変更」を選択。
   - 同様に主要都市のリストを提示し、選択させる。

コントローラの調整

```go
// internal/interfaces/controller/line_webhook.go
package controller

import (
    "encoding/json"
    "net/http"
    "projec-root/internal/usecase"
    "projec-root/pkg/lineapi"
    "fmt"
)

type LineWebhookController struct {
    AuthUseCase     *usecase.AuthUserUseCase
    NotifyUseCase   *usecase.NotifyWeatherUseCase
    SetLocationUseCase *usecase.SetLocationUseCase
    LineClient      *lineapi.Client
}

func NewLineWebhookController(authUC *usecase.AuthUserUseCase, notifyUC *usecase.NotifyWeatherUseCase, setLocUC *usecase.SetLocationUseCase, lineClient *lineapi.Client) *LineWebhookController {
    return &LineWebhookController{
        AuthUseCase:     authUC,
        NotifyUseCase:   notifyUC,
        SetLocationUseCase: setLocUC,
        LineClient:      lineClient,
    }
}

func (c *LineWebhookController) HandleWebhook(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        Events []struct {
            Type     string `json:"type"`
            Source   struct {
                UserID string `json:"userId"`
            } `json:"source"`
            Message struct {
                Type string `json:"type"`
                Text string `json:"text"`
            } `json:"message"`
            Postback struct {
                Data string `json:"data"`
            } `json:"postback"`
        } `json:"events"`
    }

    // シグネチャの検証（省略）

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    for _, event := range payload.Events {
        if event.Type == "message" && event.Message.Type == "text" {
            userID := event.Source.UserID
            message := event.Message.Text

            // 認証処理
            user, err := c.AuthUseCase.Authenticate(userID)
            if err != nil {
                // エラーハンドリング
                continue
            }

            // ユーザーからの入力に基づく処理
            if message == "都市を変更" {
                // 都市選択メニューを再提示
                cities, err := c.NotifyUseCase.GetAllCities()
                if err != nil {
                    fmt.Printf("Failed to get cities: %v\n", err)
                    continue
                }
                // リッチメニューやテンプレートメッセージで都市を提示
                c.LineClient.SendMessage(user.LineUserID, "都市を選択してください。")
                // 都市選択メニューの実装を追加
                continue
            }

            // その他のメッセージ処理
            // ...
        }

        if event.Type == "postback" && event.Postback.Type == "postback" {
            userID := event.Source.UserID
            data := event.Postback.Data // 例: "select_city:130000"

            // データのパース
            var action, cityID string
            fmt.Sscanf(data, "select_city:%s", &cityID)

            // 都市設定ユースケースの実行
            err := c.SetLocationUseCase.Execute(userID, cityID)
            if err != nil {
                fmt.Printf("Failed to set location for user %s: %v\n", userID, err)
                continue
            }

            // 確認メッセージの送信
            c.LineClient.SendMessage(user.LineUserID, "都市が更新されました。")
        }
    }

    w.WriteHeader(http.StatusOK)
}
```

7. ディレクトリ構成の再確認（更新）

``` shell
/projec-root
├── /cmd                 # アプリケーションのエントリーポイント（main.goなど）
├── /internal            # アプリケーション固有の内部ロジック
│     ├── /entity        # ドメインモデル（エンティティ）
│     │     └── user.go  # ユーザーエンティティ（例：LINEユーザー情報など）
│     ├── /usecase       # ユースケース
│     │     ├── notify_weather.go  # 天気通知のユースケース
│     │     ├── set_location.go    # 都市設定のユースケース
│     │     └── auth_user.go       # LINE認証処理のユースケース
│     ├── /interfaces    # インターフェースアダプタ（データ変換やリポジトリの実装）
│     │     ├── /controller # コントローラ（HTTPリクエストを処理する部分）
│     │     │     └── line_webhook.go # LINEからのWebhook処理
│     │     ├── /repository # データベースリポジトリ
│     │     │     └── user_repository.go # ユーザーデータの操作
│     │     └── /presenter # 出力変換、APIレスポンスの整形など
│     │           └── line_presenter.go # LINE向けのメッセージ整形
├── /pkg                 # 再利用可能なコード（汎用的なライブラリ）
│     ├── /weather       # 気象庁APIによる天気情報取得のロジック
│     ├── /cache         # Redisキャッシュのロジック
│     └── /lineapi       # LINE Messaging APIのクライアントラッパー
├── /configs             # 設定ファイル（YAML, JSONなど）
│     └── cities.json    # 主要都市リスト
├── /scripts             # デプロイやビルドなどのスクリプト
├── /build               # コンテナやCI/CDで使用するビルド関連のファイル
│     └── kubernetes     # Kubernetesマニフェストファイル
├── /test                # テストコード
│     ├── /unit          # ユニットテスト
│     └── /integration   # 統合テスト
```
