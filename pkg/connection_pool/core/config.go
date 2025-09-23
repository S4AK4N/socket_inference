package core

import (
	"time"
)

// PoolConfig は接続プールの設定
// 設定に関する全ての責務を担当
type PoolConfig struct {
	// 最大接続数
	MaxPoolSize int

	// 接続のタイムアウト
	ConnectTimeout time.Duration

	// アイドル接続の最大保持時間
	IdleTimeout time.Duration

	// 接続の最大生存時間
	MaxLifetime time.Duration

	// クリーンアップの間隔
	CleanupInterval time.Duration

	// サーバーURL
	ServerURL string
}

// Validate は設定の妥当性を検証
func (c *PoolConfig) Validate() error {
	if c.MaxPoolSize <= 0 {
		c.MaxPoolSize = 10
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = 5 * time.Minute
	}
	if c.MaxLifetime <= 0 {
		c.MaxLifetime = 30 * time.Minute
	}
	if c.CleanupInterval <= 0 {
		c.CleanupInterval = 1 * time.Minute
	}
	return nil
}

// ConnectionOptions は接続作成時の設定オプション
type ConnectionOptions struct {
	ServerURL      string
	ConnectTimeout time.Duration
	Headers        map[string][]string
	Subprotocols   []string
}

// PoolStats はプールの統計情報
type PoolStats struct {
	TotalConnections  int
	ActiveConnections int
	IdleConnections   int
	MaxPoolSize       int
	CreatedCount      int64
	ReusedCount       int64
	ErrorCount        int64
}
