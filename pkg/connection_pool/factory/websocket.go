package factory

import (
	"context"
	"fmt"
	"net/url"

	"socket_inference/pkg/connection_pool/core"

	"github.com/coder/websocket"
)

// WebSocketConnectionFactory は WebSocket 接続作成の具体実装
// WebSocket接続作成の詳細なロジックを担当
type WebSocketConnectionFactory struct {
	options core.ConnectionOptions
}

// NewWebSocketConnectionFactory は新しい WebSocket 接続ファクトリを作成
func NewWebSocketConnectionFactory(options core.ConnectionOptions) ConnectionFactory {
	return &WebSocketConnectionFactory{
		options: options,
	}
}

// CreateConnection は新しいWebSocket接続を作成
func (f *WebSocketConnectionFactory) CreateConnection(ctx context.Context) (*websocket.Conn, error) {
	u, err := url.Parse(f.options.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("URL解析エラー: %w", err)
	}

	dialCtx, cancel := context.WithTimeout(ctx, f.options.ConnectTimeout)
	defer cancel()

	conn, _, err := websocket.Dial(dialCtx, u.String(), &websocket.DialOptions{
		HTTPHeader:   f.options.Headers,
		Subprotocols: f.options.Subprotocols,
	})
	if err != nil {
		return nil, fmt.Errorf("WebSocket接続エラー: %w", err)
	}

	return conn, nil
}

// ValidateConnection は接続が有効かどうかを検証
func (f *WebSocketConnectionFactory) ValidateConnection(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("接続がnilです")
	}

	// パフォーマンス重視：検証をスキップ
	// 実際の使用時にエラーが発生した場合に接続を削除する方が効率的
	return nil

	// オプション：軽量な検証を行いたい場合
	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// defer cancel()
	// return conn.Ping(ctx)
}

// GetConnectionOptions は接続オプションを返す
func (f *WebSocketConnectionFactory) GetConnectionOptions() core.ConnectionOptions {
	return f.options
}
