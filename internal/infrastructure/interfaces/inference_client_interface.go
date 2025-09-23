package interfaces

import (
	"context"
	"socket_inference/internal/model"
)

// InferenceClient 推論サーバーとの通信インターフェース
// ViewModelがこのインターフェースに依存し、具体的な実装（gRPC等）は隠蔽
type InferenceClient interface {
	// SendInferenceRequest 推論リクエストをサーバーに送信
	SendInferenceRequest(ctx context.Context, request *model.InferenceRequest) (*model.InferenceResponse, error)

	// SendBatchInferenceRequest バッチ推論リクエストを送信
	SendBatchInferenceRequest(ctx context.Context, batch *model.AudioBatch) (*model.InferenceResponse, error)

	// Connect 推論サーバーに接続
	Connect(ctx context.Context) error

	// Disconnect 推論サーバーから切断
	Disconnect() error

	// IsConnected 接続状態を確認
	IsConnected() bool

	// GetServerStatus サーバーの状態を取得
	GetServerStatus() (string, error)
}

// ConnectionConfig 接続設定のインターフェース
type ConnectionConfig interface {
	// GetServerAddress サーバーアドレスを取得
	GetServerAddress() string

	// GetTimeout タイムアウト値を取得
	GetTimeout() int

	// GetRetryPolicy リトライポリシーを取得
	GetRetryPolicy() map[string]interface{}
}
