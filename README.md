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
git clone https://github.com/S4AK4N/socket_inference.git
cd socket_inference

# 依存関係取得
go mod tidy

# ビルド
go build -o main .
```

### 3. 基本実行
```bash
# サーバー起動
./main

# 別ターミナルでテストクライアント実行
go run cmd/test_client/main.go
```

## 🔧 パラメータチューニング

環境変数でシステム動作をカスタマイズできます：

```bash
# パラメータチューニング用クライアント実行
CLIENT_COUNT=5 \
BATCH_SIZE=15 \
FLUSH_TIMEOUT=3s \
go run cmd/tuning_client/main.go
```

### 利用可能な環境変数

| 環境変数 | デフォルト値 | 説明 |
|----------|-------------|------|
| `SERVER_PORT` | `8080` | WebSocketサーバーポート |
| `CLIENT_COUNT` | `3` | 同時接続クライアント数 |
| `BATCH_SIZE` | `10` | 音声チャンクバッチサイズ |
| `FLUSH_TIMEOUT` | `2s` | バッチフラッシュタイムアウト |
| `TEST_DURATION` | `30s` | テスト実行時間 |

## 🧪 テストスクリプト

```bash
# 壁打ちテスト（基本動作確認）
./scripts/test_wallpunch.sh

# パラメータチューニングテスト
./scripts/tune_test.sh
```

## 📊 パフォーマンス監視

システムは以下のメトリクスを自動収集します：

- **スループット**: 秒間処理チャンク数
- **レイテンシ**: リクエスト〜レスポンス時間
- **エラー率**: 失敗リクエスト割合
- **接続数**: アクティブWebSocket接続数

## 🔗 ドキュメント

詳細な技術仕様は以下をご覧ください：

- [📋 API仕様](docs/API.md) - WebSocket APIの詳細仕様
- [🏛️ アーキテクチャ設計](docs/ARCHITECTURE.md) - Clean Architecture実装詳細
- [🧪 テスト戦略](docs/TESTING.md) - テストケースとシナリオ
- [⚙️ チューニングガイド](docs/TUNING.md) - パフォーマンス最適化手順

## 🛠️ 開発環境

### 依存関係
- `github.com/coder/websocket` - WebSocket実装
- Go標準ライブラリ

### ビルド・実行
```bash
# 開発用サーバー起動
go run main.go

# プロダクションビルド  
go build -ldflags="-s -w" -o socket_inference

# テストスイート実行
go test ./...
```

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. フィーチャーブランチ作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエスト作成

## 📄 ライセンス

このプロジェクトは [MIT License](LICENSE) の下で公開されています。

## 🏷️ バージョン履歴

- **v1.0.0** - 初回リリース
  - Clean Architecture実装
  - WebSocket音声ストリーミング
  - パラメータチューニング機能
  - 並行テストクライアント

---

💡 **Tip**: 詳細な使用方法は各ドキュメントファイルをご確認ください。
