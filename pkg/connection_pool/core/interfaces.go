package core

import (
	"context"
)

// ConnectionPool は WebSocket 接続プールのメインインターフェース
// 外部から使用される公開APIを定義
type ConnectionPool interface {
	// Get は利用可能な接続を取得する
	Get(ctx context.Context) (*PooledConnection, error)

	// Put は接続をプールに戻す
	Put(conn *PooledConnection) error

	// Close は接続を完全に閉じてプールから削除
	Close(conn *PooledConnection) error

	// Cleanup は古い接続や使用されていない接続を清掃
	Cleanup() error

	// Shutdown は全ての接続を閉じてプールをシャットダウン
	Shutdown() error

	// Stats はプールの統計情報を返す
	Stats() PoolStats
}
