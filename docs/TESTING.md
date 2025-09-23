# 🎯 壁打ちテスト - 動作確認ガイド

Clean Architectureで構築された音声ストリーミングシステムの動作確認用ガイドです。

## 🏗️ アーキテクチャ概要

```
Domain (model) ← Use Case (viewmodel) ← Interface Adapter (view) ← Infrastructure
```

### 主要コンポーネント:
- **Domain Layer (model/)**: AudioClient, AudioBatch, InferenceRequest/Response
- **Use Case Layer (viewmodel/)**: ビジネスロジック（バッチ処理、推論管理）
- **Interface Adapter Layer (view/)**: WebSocketハンドラー、HTTPサーバー
- **Infrastructure Layer (infrastructure/)**: gRPCクライアント（プレースホルダー）

## 🚀 動作確認手順

### 1. サーバー起動
```bash
# メインサーバーを起動
go run main.go
```

### 2. テストクライアント実行
```bash
# 別ターミナルで壁打ちテスト実行
go run cmd/test_client/main.go
```

## 📋 確認ポイント

### ✅ 正常動作確認項目
1. **WebSocket接続**: クライアントがサーバーに正常接続
2. **音声データ受信**: 10チャンクの模擬音声データ送信
3. **バッチ処理**: 10チャンク到達でバッチ生成
4. **推論処理**: 模擬推論サーバーへの送信ログ
5. **タイムアウト処理**: 2秒タイムアウトでの自動フラッシュ

### 📊 期待されるログ出力
```
クライアント test-client-001 が接続しました
バッチ準備完了: クライアント=test-client-001, チャンク数=10
推論サーバーにバッチ送信: クライアント=test-client-001, チャンク数=10
推論処理を開始しました
```

## 🔧 デバッグのヒント

### エラーが出た場合：
1. **ポート8080が使用中**: 他のプロセスを終了するか、port番号変更
2. **依存関係エラー**: `go mod tidy`を実行
3. **インポートエラー**: Clean Architecture構造を確認

### ログが少ない場合：
- バッチサイズ10または2秒タイムアウトで処理される
- クライアント数を増やして並行処理をテスト

## 🎯 拡張テスト

### 複数クライアントテスト:
```bash
# 複数のテストクライアントを同時実行
for i in {1..3}; do
    go run cmd/test_client/main.go &
done
wait
```

### カスタムテストデータ:
`cmd/test_client/main.go`の`audioData`配列を変更してテスト

## 📚 コード理解のための順序

1. **main.go**: エントリーポイントと依存注入
2. **model/**: ドメインエンティティ
3. **viewmodel/coordinator/**: 全体調整ロジック
4. **view/handlers/websocket/**: WebSocket処理
5. **infrastructure/grpc/**: 外部システム連携