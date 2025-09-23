package connection_pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"socket_inference/pkg/connection_pool/core"
	"socket_inference/pkg/connection_pool/factory"
	"socket_inference/pkg/connection_pool/manager"

	"github.com/google/uuid"
)

// ConnectionPool は待機機能付きの接続プール実装
type ConnectionPool struct {
	config    core.PoolConfig
	factory   factory.ConnectionFactory
	manager   manager.PoolManager
	waitQueue chan chan *core.PooledConnection // 待機キュー
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewConnectionPool は新しい接続プールを作成
func NewConnectionPool(config core.PoolConfig) *ConnectionPool {
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("設定が無効です: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 接続オプション設定
	connectionOptions := core.ConnectionOptions{
		ServerURL:      config.ServerURL,
		ConnectTimeout: config.ConnectTimeout,
		Headers: map[string][]string{
			"User-Agent": {"connection-pool-client/1.0"},
		},
	}

	pool := &ConnectionPool{
		config:    config,
		factory:   factory.NewWebSocketConnectionFactory(connectionOptions),
		manager:   manager.NewDefaultPoolManager(config.MaxPoolSize),
		waitQueue: make(chan chan *core.PooledConnection, 100), // 無制限に近い待機キューサイズ
		ctx:       ctx,
		cancel:    cancel,
	}

	return pool
}

// Get は利用可能な接続を取得（待機機能付き）
func (p *ConnectionPool) Get(ctx context.Context) (*core.PooledConnection, error) {
	// まず即座に利用可能な接続を探す
	if conn := p.tryGetConnection(); conn != nil {
		return conn, nil
	}

	// 新しい接続を作成できるかチェック
	if !p.manager.IsPoolFull() {
		return p.createNewConnection(ctx)
	}

	// プールが満杯の場合、待機する
	return p.waitForConnection(ctx)
}

// tryGetConnection は即座に利用可能な接続を探す
func (p *ConnectionPool) tryGetConnection() *core.PooledConnection {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn := p.manager.FindAvailableConnection(); conn != nil {
		conn.MarkAsUsed()
		return conn
	}
	return nil
}

// createNewConnection は新しい接続を作成
func (p *ConnectionPool) createNewConnection(ctx context.Context) (*core.PooledConnection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := p.factory.CreateConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("接続作成失敗: %w", err)
	}

	pooledConn := &core.PooledConnection{
		ID:        uuid.New().String(),
		Conn:      conn,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		InUse:     true,
	}

	if err := p.manager.AddConnection(pooledConn); err != nil {
		conn.Close(1000, "プール追加失敗")
		return nil, fmt.Errorf("プール追加失敗: %w", err)
	}

	return pooledConn, nil
}

// waitForConnection は接続が利用可能になるまで待機
func (p *ConnectionPool) waitForConnection(ctx context.Context) (*core.PooledConnection, error) {
	// 待機用チャネルを作成
	connChan := make(chan *core.PooledConnection, 1)

	// 10秒のタイムアウトを設定
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 待機キューに追加
	select {
	case p.waitQueue <- connChan:
		// キューに追加成功
	case <-waitCtx.Done():
		return nil, fmt.Errorf("待機キューへの追加がタイムアウトしました")
	}

	// 接続が利用可能になるまで待機
	select {
	case conn := <-connChan:
		return conn, nil
	case <-waitCtx.Done():
		return nil, fmt.Errorf("接続待機がタイムアウトしました（10秒）")
	}
}

// Put は接続をプールに戻す
func (p *ConnectionPool) Put(conn *core.PooledConnection) error {
	if conn == nil {
		return fmt.Errorf("接続がnilです")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 接続の有効性を確認
	if err := p.factory.ValidateConnection(conn.Conn); err != nil {
		p.manager.RemoveConnection(conn.ID)
		if conn.Conn != nil {
			conn.Conn.Close(1000, "接続無効")
		}
		return fmt.Errorf("接続が無効です: %w", err)
	}

	conn.MarkAsIdle()

	// 待機中のクライアントがいるかチェック
	select {
	case connChan := <-p.waitQueue:
		// 待機中のクライアントに接続を渡す
		conn.MarkAsUsed()
		connChan <- conn
	default:
		// 待機中のクライアントがいない場合はそのままプールに戻す
	}

	return nil
}

// Close は接続を完全に閉じてプールから削除
func (p *ConnectionPool) Close(conn *core.PooledConnection) error {
	if conn == nil {
		return fmt.Errorf("接続がnilです")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.manager.RemoveConnection(conn.ID); err != nil {
		return fmt.Errorf("プールからの削除失敗: %w", err)
	}

	if conn.Conn != nil {
		conn.Conn.Close(1000, "プールから削除")
	}

	return nil
}

// Cleanup は古い接続や無効な接続を清掃
func (p *ConnectionPool) Cleanup() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	connections := p.manager.GetAllConnections()
	var toRemove []string

	for _, conn := range connections {
		if conn.IsExpired(p.config.MaxLifetime) || conn.IsIdle(p.config.IdleTimeout) {
			toRemove = append(toRemove, conn.ID)
			if conn.Conn != nil {
				conn.Conn.Close(1000, "クリーンアップ")
			}
		}
	}

	for _, connID := range toRemove {
		p.manager.RemoveConnection(connID)
	}

	return nil
}

// Shutdown は全ての接続を閉じてプールをシャットダウン
func (p *ConnectionPool) Shutdown() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cancel()

	// 全ての接続を閉じる
	connections := p.manager.GetAllConnections()
	for _, conn := range connections {
		if conn.Conn != nil {
			conn.Conn.Close(1000, "プールシャットダウン")
		}
	}

	// 待機キューをクリア
	close(p.waitQueue)

	return nil
}

// Stats はプールの統計情報を返す
func (p *ConnectionPool) Stats() core.PoolStats {
	connections := p.manager.GetAllConnections()
	total := len(connections)
	active := 0
	idle := 0

	for _, conn := range connections {
		if conn.InUse {
			active++
		} else {
			idle++
		}
	}

	return core.PoolStats{
		TotalConnections:  total,
		ActiveConnections: active,
		IdleConnections:   idle,
		MaxPoolSize:       p.config.MaxPoolSize,
		// 統計カウンターは別途実装が必要
	}
}
