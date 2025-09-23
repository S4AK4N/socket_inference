package manager

import (
	"socket_inference/pkg/connection_pool/core"
)

// PoolManager はプールの管理機能を抽象化するインターフェース
// 接続の追加・削除・検索などプール管理の責務を担当
type PoolManager interface {
	// AddConnection は接続をプールに追加
	AddConnection(conn *core.PooledConnection) error

	// RemoveConnection は接続をプールから削除
	RemoveConnection(connID string) error

	// FindAvailableConnection は利用可能な接続を検索
	FindAvailableConnection() *core.PooledConnection

	// GetAllConnections は全ての接続を返す
	GetAllConnections() []*core.PooledConnection

	// GetConnectionCount は現在の接続数を返す
	GetConnectionCount() int

	// IsPoolFull はプールが満杯かどうかを確認
	IsPoolFull() bool
}
