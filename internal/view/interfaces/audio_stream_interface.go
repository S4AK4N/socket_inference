package interfaces

import "socket_inference/internal/model"

// AudioStreamHandler 音声ストリーミング処理の共通インターフェース
// WebSocket、gRPC等の異なるプロトコルで共通利用可能
type AudioStreamHandler interface {
	// HandleConnection 接続を処理し、音声ストリーミングを開始
	HandleConnection(connectionData interface{}) error
}

// AudioViewModelInterface 音声ViewModelのインターフェースを定義
type AudioViewModelInterface interface {
	RegisterClient(client *model.AudioClient)
	UnregisterClient(client *model.AudioClient)
	ProcessAudioData(clientID string, audioData []byte)
}
