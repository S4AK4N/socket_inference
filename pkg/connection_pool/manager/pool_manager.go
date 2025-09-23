package manager

import (
	"fmt"
	"sync"

	"socket_inference/pkg/connection_pool/core"
)

// DefaultPoolManager はプール管理の標準実装
// スレッドセーフな接続プール管理を提供
type DefaultPoolManager struct {
	connections []*core.PooledConnection
	maxSize     int
	mu          sync.RWMutex
}

// NewDefaultPoolManager は新しいプールマネージャーを作成
func NewDefaultPoolManager(maxSize int) PoolManager {
	return &DefaultPoolManager{
		connections: make([]*core.PooledConnection, 0, maxSize),
		maxSize:     maxSize,
	}
}

// AddConnection は接続をプールに追加
func (pm *DefaultPoolManager) AddConnection(conn *core.PooledConnection) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.connections) >= pm.maxSize {
		return fmt.Errorf("プールが満杯です (最大: %d)", pm.maxSize)
	}

	pm.connections = append(pm.connections, conn)
	return nil
}

// RemoveConnection は接続をプールから削除
func (pm *DefaultPoolManager) RemoveConnection(connID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, conn := range pm.connections {
		if conn.ID == connID {
			// スライスから削除
			pm.connections = append(pm.connections[:i], pm.connections[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("接続ID %s が見つかりません", connID)
}

// FindAvailableConnection は利用可能な接続を検索
func (pm *DefaultPoolManager) FindAvailableConnection() *core.PooledConnection {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, conn := range pm.connections {
		if !conn.InUse {
			return conn
		}
	}

	return nil
}

// GetAllConnections は全ての接続を返す
func (pm *DefaultPoolManager) GetAllConnections() []*core.PooledConnection {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// コピーを返してスレッドセーフにする
	result := make([]*core.PooledConnection, len(pm.connections))
	copy(result, pm.connections)
	return result
}

// GetConnectionCount は現在の接続数を返す
func (pm *DefaultPoolManager) GetConnectionCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.connections)
}

// IsPoolFull はプールが満杯かどうかを確認
func (pm *DefaultPoolManager) IsPoolFull() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.connections) >= pm.maxSize
}
