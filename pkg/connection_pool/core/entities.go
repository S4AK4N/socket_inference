package core

import (
	"time"

	"github.com/coder/websocket"
)

// PooledConnection は接続プールで管理される接続のエンティティ
// ドメインの中核となるエンティティで、接続の状態と生存期間を管理
type PooledConnection struct {
	// 接続の一意識別子
	ID string

	// 実際のWebSocket接続
	Conn *websocket.Conn

	// 接続が作成された時刻
	CreatedAt time.Time

	// 最後に使用された時刻
	LastUsed time.Time

	// 接続が現在使用中かどうか
	InUse bool
}

// IsExpired は接続が指定された最大生存時間を超えているかチェック
func (pc *PooledConnection) IsExpired(maxLifetime time.Duration) bool {
	return time.Since(pc.CreatedAt) > maxLifetime
}

// IsIdle は接続が指定されたアイドル時間を超えているかチェック
func (pc *PooledConnection) IsIdle(idleTimeout time.Duration) bool {
	return !pc.InUse && time.Since(pc.LastUsed) > idleTimeout
}

// MarkAsUsed は接続を使用中としてマーク
func (pc *PooledConnection) MarkAsUsed() {
	pc.InUse = true
	pc.LastUsed = time.Now()
}

// MarkAsIdle は接続をアイドル状態としてマーク
func (pc *PooledConnection) MarkAsIdle() {
	pc.InUse = false
	pc.LastUsed = time.Now()
}
