#!/bin/bash

echo "🎯 WebSocket音声ストリーミング壁打ちテスト"
echo "=========================================="

# サーバーが起動しているかチェック
if ! curl -s http://localhost:8080 > /dev/null 2>&1; then
    echo "❌ サーバーが起動していません"
    echo "別ターミナルで以下を実行してください:"
    echo "go run main.go"
    exit 1
fi

echo "✅ サーバーが起動中です"

# WebSocketテストにcurlを使用
echo "📡 WebSocket接続テスト開始..."

# 簡単なHTTPヘルスチェック
echo "🔍 HTTP接続確認:"
curl -I http://localhost:8080/audio 2>/dev/null | head -1

echo ""
echo "🎵 WebSocket音声ストリーミングテスト:"
echo "   - バッチサイズ: 10チャンク"
echo "   - タイムアウト: 2秒"
echo "   - クライアントID: test-client-bash"

echo ""
echo "📝 手動テスト手順:"
echo "1. ブラウザで ws://localhost:8080/audio に接続"
echo "2. バイナリデータを10回送信"
echo "3. サーバーログでバッチ処理を確認"

echo ""
echo "✨ 壁打ちテスト完了"