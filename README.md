# 🎵 Socket Inference - リアルタイム音声推論システム

Clean Architectureに基づいて構築されたWebSocket音声ストリーミング推論システムです。

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-green.svg)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![WebSocket](https://img.shields.io/badge/Protocol-WebSocket-orange.svg)](https://tools.ietf.org/html/rfc6455)

## 🏗️ アーキテクチャ概要

```
Domain (model) ← Use Case (viewmodel) ← Interface Adapter (view) ← Infrastructure
```

### 📁 プロジェクト構造

```
socket_inference/
├── cmd/                    # アプリケーションエントリーポイント
│   ├── test_client/       # 単純テストクライアント
│   └── tuning_client/     # パラメータチューニング用クライアント
├── internal/              # プライベートパッケージ
│   ├── model/            # Domain Layer (最抽象)
│   ├── viewmodel/        # Use Case Layer
│   ├── view/             # Interface Adapter Layer  
│   ├── infrastructure/   # Infrastructure Layer (最具象)
│   └── config/           # 設定管理
├── docs/                 # ドキュメント
├── scripts/              # 実行スクリプト
└── README.md
```

## 🚀 クイックスタート

### 1. 前提条件
- Go 1.25+ インストール済み
- Git インストール済み

### 2. セットアップ
```bash
# リポジトリクローン
git clone https://github.com/yoshidarimare/socket_inference.git
cd socket_inference

# 依存関係取得
go mod tidy

# ビルド
go build -o main .
```

### 3. サーバー起動
```bash
# デフォルト設定で起動
./main

# 環境変数でカスタマイズ
BATCH_SIZE=20 FLUSH_TIMEOUT=1s ./main
```

### 4. テスト実行
```bash
# 基本的な動作確認
./scripts/test_wallpunch.sh

# パラメータチューニングテスト
CLIENT_COUNT=5 ./scripts/tune_test.sh
```

## 📊 機能

### ✅ 実装済み機能
- **WebSocket音声ストリーミング**: リアルタイム音声データ受信
- **音声バッチ処理**: 設定可能なバッチサイズとタイムアウト
- **並行クライアント対応**: 複数クライアント同時接続
- **推論処理パイプライン**: gRPC推論サーバー連携（プレースホルダー）
- **パラメータチューニング**: 環境変数による動的設定
- **パフォーマンス監視**: リアルタイム統計とメトリクス

### 🔄 処理フロー
1. クライアントがWebSocket接続 (`:8080/audio`)
2. 音声データをバイナリ形式でストリーミング送信
3. サーバーがデータをバッチ化（10チャンクまたは2秒タイムアウト）
4. バッチを推論処理パイプラインに送信
5. 結果をクライアントに返却（将来実装）

## ⚙️ 設定

### 環境変数
| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SERVER_PORT` | `8080` | サーバーポート |
| `BATCH_SIZE` | `10` | 音声バッチサイズ |
| `FLUSH_TIMEOUT` | `2s` | バッチフラッシュタイムアウト |
| `MAX_CLIENTS` | `100` | 最大同時接続数 |
| `GRPC_SERVER` | `localhost:50051` | gRPCサーバーアドレス |

詳細な設定については [docs/TUNING.md](docs/TUNING.md) を参照してください。

## 🧪 テストとデバッグ

### 基本テスト
```bash
# サーバー起動
go run main.go

# 別ターミナルで動作確認
go run cmd/test_client/main.go
```

### パフォーマンステスト
```bash
# 高負荷テスト
CLIENT_COUNT=20 CHUNK_INTERVAL=10ms ./scripts/tune_test.sh

# スループットテスト  
CHUNK_SIZE=8192 TEST_DURATION=60s ./scripts/tune_test.sh
```

## 📚 ドキュメント

- [🏗️ アーキテクチャ設計](docs/ARCHITECTURE.md) - Clean Architecture実装詳細
- [🔧 パラメータチューニング](docs/TUNING.md) - 性能最適化ガイド
- [🧪 テストガイド](docs/TESTING.md) - テスト実行方法
- [📡 API仕様](docs/API.md) - WebSocket API仕様

## 🔧 開発

### 依存関係管理
```bash
# 依存関係更新
go mod tidy

# セキュリティ検査
go mod verify
```

### コード品質
```bash
# フォーマット
go fmt ./...

# Lint検査
golangci-lint run

# テスト実行
go test ./...
```

## 🎯 今後の開発予定

- [ ] 実際のgRPC推論サーバー実装
- [ ] 音声形式対応拡張（WAV, MP3等）
- [ ] 認証・認可機能
- [ ] ヘルスチェックエンドポイント
- [ ] Docker化対応
- [ ] Kubernetes対応
- [ ] モニタリング・ログ改善

## 🤝 貢献

プルリクエストやイシューは歓迎です！

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🙏 謝辞

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [WebSocket library](https://github.com/coder/websocket) by Coder

---

**Built with ❤️ using Go and Clean Architecture principles**