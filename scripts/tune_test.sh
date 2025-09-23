#!/bin/bash

echo "🎯 パラメータチューニング用テストスイート"
echo "========================================"

# デフォルト値
SERVER_URL="${SERVER_URL:-ws://localhost:8080/audio}"
CLIENT_COUNT="${CLIENT_COUNT:-3}"
CHUNKS_PER_CLIENT="${CHUNKS_PER_CLIENT:-15}"
CHUNK_INTERVAL="${CHUNK_INTERVAL:-100ms}"
CHUNK_SIZE="${CHUNK_SIZE:-1024}"
TEST_DURATION="${TEST_DURATION:-10s}"

echo "📋 現在の設定:"
echo "   SERVER_URL=$SERVER_URL"
echo "   CLIENT_COUNT=$CLIENT_COUNT"
echo "   CHUNKS_PER_CLIENT=$CHUNKS_PER_CLIENT"
echo "   CHUNK_INTERVAL=$CHUNK_INTERVAL"
echo "   CHUNK_SIZE=$CHUNK_SIZE"
echo "   TEST_DURATION=$TEST_DURATION"
echo ""

# サーバーが起動しているかチェック
echo "🔍 サーバー状態確認..."
if ! curl -s http://localhost:8080 > /dev/null 2>&1; then
    echo "❌ サーバーが起動していません"
    echo "別ターミナルで以下を実行してください:"
    echo "go run main.go"
    exit 1
fi
echo "✅ サーバー起動確認"

echo ""
echo "🚀 パラメータチューニングテスト開始..."

# 環境変数を設定してクライアント実行
export SERVER_URL CLIENT_COUNT CHUNKS_PER_CLIENT CHUNK_INTERVAL CHUNK_SIZE TEST_DURATION
go run cmd/tuning_client/main.go

echo ""
echo "💡 パラメータ調整例:"
echo "   CLIENT_COUNT=5 CHUNK_INTERVAL=50ms ./tune_test.sh    # 高負荷テスト"
echo "   CLIENT_COUNT=1 CHUNKS_PER_CLIENT=50 ./tune_test.sh   # 長時間テスト"
echo "   CHUNK_SIZE=2048 TEST_DURATION=30s ./tune_test.sh     # 大容量テスト"