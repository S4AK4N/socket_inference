package interfaces

import "socket_inference/internal/model"

// ClientManager クライアント接続管理のインターフェース
type ClientManager interface {
	// RegisterClient 新しいクライアントを登録
	RegisterClient(client *model.AudioClient)

	// UnregisterClient クライアントの登録を解除
	UnregisterClient(client *model.AudioClient)

	// GetConnectedClients 接続中のクライアント一覧を取得
	GetConnectedClients() []*model.AudioClient

	// GetClientCount 接続中のクライアント数を取得
	GetClientCount() int
}
