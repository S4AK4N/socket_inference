package factory

import (
	"context"

	"socket_inference/pkg/connection_pool/core"

	"github.com/coder/websocket"
)

// ConnectionFactory はWebSocket接続の作成を抽象化するインターフェース
// 接続作成に関する全ての責務を担当
type ConnectionFactory interface {
	// CreateConnection は新しいWebSocket接続を作成
	CreateConnection(ctx context.Context) (*websocket.Conn, error)

	// ValidateConnection は接続が有効かどうかを検証
	ValidateConnection(conn *websocket.Conn) error

	// GetConnectionOptions は接続オプションを返す
	GetConnectionOptions() core.ConnectionOptions
}
