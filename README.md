## プロジェクト概要

このプロジェクトは、ユーザーが指定した時刻と地域に基づいて天気情報をLINEを通じて通知するLINE BOTの開発を目的としている。LINEアカウントを利用した認証機能を活用し、シンプルで効率的な天気情報通知システムを構築する。

## 機能概要

- **天気通知機能**: ユーザーが指定した時刻および地域に基づき、天気情報をLINEで通知する。
- **ユーザー認証**: LINEアカウント情報を利用したユーザー認証。
- **通知スケジュール管理**: ユーザーごとに通知の時刻と地域を設定し、指定した時間に天気情報を送信。
- **位置情報設定機能**: ユーザーが位置情報を登録し、天気情報を通知するための地域を設定。

## 技術スタック

### 1. アプリケーションアーキテクチャ

- **コンテナ技術**: Kubernetesを利用してコンテナを管理。HelmやKustomizeを使用してKubernetesのマニフェストを効率的に管理し、後々の運用も見据えた設定を行う。
- **API連携**: 以下のAPIを活用して天気情報と位置情報を取得する。
  - 天気予報: OpenWeatherMap API
  - 緯度経度取得用: Google Geocoding API

### 2. システムインフラ

- **CI/CD**:
  - **CI**: GitHub Actionsにより自動テストとデプロイを効率化。
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

## 詳細設計

### 1. 天気通知機能

#### 機能概要
ユーザーが指定した時刻および地域に基づき、天気情報をLINEで通知する機能です。基本的には雨の時のみ通知します。

#### 実装のポイント

##### データ取得
- **OpenWeatherMap APIの活用**:
  - 現在の天気情報や週間予報を取得するために、適切なエンドポイントを選択します。
  - APIのレスポンスデータを解析し、必要な情報（気温、降水確率など）を抽出します。

- **APIの効率的な利用**:
  - ユーザー数が増加した場合のAPI呼び出し回数を最適化するために、キャッシュ（例：Redis）を活用します。
  - オンデマンド取得を採用し、必要なときにのみAPIを呼び出す方式とします。
  - 定期的にバッチ処理を行い、必要なデータを事前に取得・保存する方法も検討します。

##### メッセージフォーマット
- **LINE Messaging APIの利用**:
  - リッチメニューやテンプレートメッセージを活用して、視覚的にわかりやすい天気情報を提供します。
  - メッセージの内容やレイアウトをユーザーの設定に応じてカスタマイズ可能にします。

##### レート制限対策
- **APIのレート制限**:
  - OpenWeatherMap APIの無料プランではリクエスト回数に制限があるため、必要に応じて有料プランへのアップグレードを検討します。
  - キャッシュの有効活用とバックオフ戦略を実装し、レート制限に対処します。

##### スケーリングとパフォーマンス
- **非同期処理の導入**:
  - 通知の送信を非同期に処理することで、システムの応答性を向上させます。

#### 実装における課題
- **APIのレート制限**:
  - OpenWeatherMap APIの制限に注意し、必要に応じて有料プランを検討。
  - ただお金は使いたくない。
  
- **データの正確性と更新頻度**:
  - 天気情報の更新頻度に応じて、適切なデータ取得タイミングを設定します。
  - 目安は1時間程度。

### 2. ユーザー認証

#### 機能概要
LINEアカウント情報を利用してユーザー認証を行い、個別の設定や通知を管理します。

#### 実装のポイント

##### 認証フロー
- **LINE Loginの利用**:
  - LINEのOAuth 2.0を使用してユーザー認証を実装します。
  - ユーザーがLINEアカウントを介してログインし、必要な権限を取得します。

- **アクセストークンの管理**:
  - LINEから取得したアクセストークンを安全に保存・管理します。
  - トークンの有効期限を監視し、必要に応じてリフレッシュトークンを使用して更新します。

##### ユーザーデータの管理
- **データベース設計**:
  - ユーザー情報（LINEのユーザーID、アクセストークン、設定情報など）をPostgreSQLに保存します。
  - セキュリティを考慮し、個人情報は暗号化して保存します。

- **プライバシーとセキュリティ**:
  - GDPRなどのプライバシー規制に準拠するための措置を講じます。
  - データベースへのアクセス制御や暗号化を徹底します。

#### 実装における課題
- **認証の信頼性**:
  - 認証フロー中のエラーハンドリングを適切に行い、ユーザーにわかりやすいエラーメッセージを提供します。
  
- **スケーラビリティ**:
  - 多数のユーザーが同時に認証を行った場合でも、システムが安定して動作するように設計します。

### 3. 通知スケジュール管理

#### 機能概要
ユーザーごとに通知の時刻と地域を設定し、指定した時間に天気情報を送信します。

#### 実装のポイント

##### スケジューリングシステム
- **ジョブスケジューラーの選定**:
  - **Kubernetes jobs.batch**: 通常のJobとして登録し、APIを通じてジョブを呼び出すパターンを採用します。

- **スケジュールの動的管理**:
  - ユーザーが通知時間を変更した際に、スケジュールを動的に更新できる仕組みを構築します。

##### データ管理
- **通知設定の保存**:
  - ユーザーごとの通知時刻、地域、通知形式などの設定をデータベースに保存します。
  - スケジュールジョブをトリガーするための情報も含めます。

##### タイムゾーンの考慮
- 日本国内のみを対象とする場合、すべてのユーザーのタイムゾーンを日本標準時（JST）に統一します。
- 将来的に多地域対応を検討する場合、ユーザーごとのタイムゾーン情報を設定に含めます。

#### 実装における課題
- **時間の正確性**:
  - サーバーのタイムゾーンとユーザーのタイムゾーンの違いを適切に管理し、通知が正確な時間に送信されるようにします。

- **スケーラビリティとパフォーマンス**:
  - 大規模なユーザー数に対応するため、スケジューリングシステムのパフォーマンスを最適化します。

### 4. 位置情報設定機能

#### 機能概要
ユーザーが位置情報を登録し、天気情報を通知するための地域を設定します。

#### 実装のポイント

##### 位置情報の入力方法
- **ユーザーインターフェース**:
  - LINEのメッセージやリッチメニューを活用して、ユーザーが簡単に位置情報を入力・変更できるようにします。
  - 位置情報をテキスト入力、座標入力、または地図から選択するオプションを提供します。

- **自動検出**:
  - LINEの位置情報機能を活用して、ユーザーの現在地を自動的に取得するオプションを提供します。

##### 位置情報の処理
- **Google Geocoding APIの活用**:
  - ユーザーが入力した住所や地名を緯度経度に変換します。
  - 入力のバリデーションとエラーハンドリングを実装します。

- **データの保存と管理**:
  - ユーザーごとの位置情報（緯度経度）をデータベースに保存します。
  - 位置情報の更新履歴を管理する場合、適切なデータモデルを設計します。

##### プライバシーとセキュリティ
- **データの保護**:
  - 位置情報は個人情報とみなされるため、適切な暗号化とアクセス制御を実施します。
  
- **ユーザーの同意**:
  - 位置情報の利用について、ユーザーから明確な同意を得るプロセスを設けます。

#### 実装における課題
- **位置情報の正確性**:
  - ユーザーが入力する位置情報の正確性を確保するため、入力方法の工夫やデータの検証を行います。

- **入力の多様性**:
  - ユーザーが様々な形式で位置情報を入力する可能性があるため、柔軟な入力処理を実装します。

### 5. キャッシュ戦略

#### キャッシュのタイミング
- **オンデマンド取得**:
  - 必要なときにのみAPIを呼び出し、取得したデータをキャッシュします。
  - キャッシュが存在しない場合や有効期限が切れている場合にのみAPIを呼び出します。

#### キャッシュの有効期限
- 天気情報の性質上、データの有効期限を設定します。例えば、30分ごとにキャッシュを更新するなど、データの鮮度とキャッシュの有効期限をバランスさせます。

#### データの分割
- **地域ごとにキャッシュを分ける**:
  - 特定の地域へのリクエストが集中しても効率的にキャッシュを利用できるようにします。

### 6. スケジューリングシステム

#### 実装の流れ
1. ユーザーが通知設定を行う。
2. 通知設定に基づき、通知時刻にジョブを登録。
3. ジョブが実行され、天気情報を取得・通知。
4. 必要に応じてジョブの再登録や更新を行う。
