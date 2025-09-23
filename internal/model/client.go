package model

import "github.com/coder/websocket"

// AudioClient 音声ストリーミング用のWebSocket接続を表現
// クライアント接続とセッション管理を担当するドメインモデル
type AudioClient struct {
	Conn     *websocket.Conn // WebSocket接続
	ClientID string          // クライアント識別用ID
}
