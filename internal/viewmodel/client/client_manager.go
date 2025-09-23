package client

import (
	"log"
	"sync"

	"socket_inference/internal/model"
	"socket_inference/internal/viewmodel/interfaces"
)

// Manager クライアント接続管理の実装
type Manager struct {
	mu      sync.RWMutex
	clients map[*model.AudioClient]bool
}

// NewManager 新しいクライアントマネージャーを作成
func NewManager() interfaces.ClientManager {
	return &Manager{
		clients: make(map[*model.AudioClient]bool),
	}
}

// RegisterClient 新しいクライアントを登録
func (cm *Manager) RegisterClient(client *model.AudioClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.clients[client] = true
	log.Printf("音声クライアント接続: %s (合計: %d)", client.ClientID, len(cm.clients))
}

// UnregisterClient クライアントの登録を解除
func (cm *Manager) UnregisterClient(client *model.AudioClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.clients[client]; ok {
		delete(cm.clients, client)
		log.Printf("音声クライアント切断: %s (合計: %d)", client.ClientID, len(cm.clients))
	}
}

// GetConnectedClients 接続中のクライアント一覧を取得
func (cm *Manager) GetConnectedClients() []*model.AudioClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clients := make([]*model.AudioClient, 0, len(cm.clients))
	for client := range cm.clients {
		clients = append(clients, client)
	}
	return clients
}

// GetClientCount 接続中のクライアント数を取得
func (cm *Manager) GetClientCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.clients)
}
