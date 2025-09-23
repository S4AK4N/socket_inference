# 🎯 パラメータチューニングガイド

環境変数を使用したパラメータチューニングとパフォーマンステストのガイドです。

## 🚀 基本的な使用方法

### 1. サーバー起動（環境変数でチューニング）
```bash
# デフォルト設定でサーバー起動
go run main.go

# カスタム設定でサーバー起動
BATCH_SIZE=20 FLUSH_TIMEOUT=1s MAX_CLIENTS=200 go run main.go
```

### 2. パラメータチューニングテスト実行
```bash
# 実行権限付与
chmod +x tune_test.sh

# デフォルト設定でテスト
./tune_test.sh

# カスタム設定でテスト
CLIENT_COUNT=10 CHUNK_INTERVAL=50ms ./tune_test.sh
```

## ⚙️ 設定パラメータ

### サーバー側環境変数
| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SERVER_PORT` | `8080` | サーバーリスニングポート |
| `BATCH_SIZE` | `10` | 音声バッチサイズ（チャンク数） |
| `FLUSH_TIMEOUT` | `2s` | バッチフラッシュタイムアウト |
| `MAX_CLIENTS` | `100` | 最大同時接続クライアント数 |
| `BUFFER_SIZE` | `100` | チャネルバッファサイズ |
| `GRPC_SERVER` | `localhost:50051` | gRPCサーバーアドレス |
| `GRPC_TIMEOUT` | `30s` | gRPCタイムアウト |

### クライアント側環境変数
| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SERVER_URL` | `ws://localhost:8080/audio` | WebSocketサーバーURL |
| `CLIENT_COUNT` | `3` | 並行実行クライアント数 |
| `CHUNKS_PER_CLIENT` | `15` | クライアント毎の送信チャンク数 |
| `CHUNK_INTERVAL` | `100ms` | チャンク送信間隔 |
| `CHUNK_SIZE` | `1024` | チャンクサイズ（バイト） |
| `TEST_DURATION` | `10s` | テスト継続時間 |

## 🧪 チューニングシナリオ例

### 高負荷テスト
```bash
# 多数クライアント + 高頻度送信
CLIENT_COUNT=20 CHUNK_INTERVAL=10ms TEST_DURATION=30s ./tune_test.sh
```

### スループットテスト
```bash
# 大容量データ + 長時間
CHUNK_SIZE=8192 CHUNKS_PER_CLIENT=100 TEST_DURATION=60s ./tune_test.sh
```

### バッチ処理最適化テスト
```bash
# サーバー側バッチサイズ調整
BATCH_SIZE=50 FLUSH_TIMEOUT=5s go run main.go &
sleep 2
CLIENT_COUNT=5 CHUNKS_PER_CLIENT=60 ./tune_test.sh
```

### 低レイテンシテスト
```bash
# 小さなバッチ + 短いタイムアウト
BATCH_SIZE=5 FLUSH_TIMEOUT=500ms go run main.go &
sleep 2
CLIENT_COUNT=3 CHUNK_INTERVAL=50ms ./tune_test.sh
```

### 耐久性テスト
```bash
# 長時間稼働テスト
TEST_DURATION=300s CHUNKS_PER_CLIENT=1000 ./tune_test.sh
```

## 📊 パフォーマンス指標

### 監視すべきメトリクス
- **スループット**: KB/s でのデータ転送率
- **レイテンシ**: バッチ処理完了までの時間
- **エラー率**: 送信失敗 / 総送信数
- **同時接続数**: 並行クライアント数
- **リソース使用率**: CPU/メモリ使用量

### ログ出力例
```
📊 [Client-001] 統計: チャンク=15, バイト=15360, エラー=0, 期間=1.52s, スループット=9.89KB/s
🏆 全体統計:
   - 総チャンク数: 45
   - 総バイト数: 46080 (0.04 MB)
   - 総エラー数: 0
   - 実行時間: 1.61秒
   - 全体スループット: 28.04 KB/s
   - 平均クライアントスループット: 9.35 KB/s
```

## 🔧 最適化のヒント

### 1. バッチサイズ最適化
- **小さいバッチ**: 低レイテンシ、高CPU使用率
- **大きいバッチ**: 高スループット、メモリ使用量増加

### 2. タイムアウト調整
- **短いタイムアウト**: リアルタイム性向上
- **長いタイムアウト**: バッチ効率向上

### 3. 並行性調整
- **クライアント数**: ネットワーク帯域とサーバー処理能力のバランス
- **バッファサイズ**: メモリ使用量と処理効率のトレードオフ

## 🚀 継続的パフォーマンステスト

### CI/CDパイプライン組み込み例
```bash
#!/bin/bash
# performance_test.sh

# ベースラインテスト
echo "ベースラインテスト実行中..."
./tune_test.sh > baseline_result.txt

# 負荷テスト
echo "負荷テスト実行中..."
CLIENT_COUNT=10 ./tune_test.sh > load_result.txt

# 結果比較とアラート
# (結果解析スクリプトの実行)
```

このような設定により、開発フェーズに応じた柔軟なパラメータチューニングが可能になります。