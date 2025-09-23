package interfaces

import (
	"context"
	"socket_inference/internal/model"
)

// AudioProcessor 音声データ処理のインターフェース
type AudioProcessor interface {
	// ProcessAudioData 音声データを処理してバッチ化
	ProcessAudioData(clientID string, audioData []byte)

	// GetBatchReady 完成したバッチを受信するチャネルを取得
	GetBatchReady() <-chan *model.AudioBatch

	// StartProcessing バックグラウンド処理を開始
	StartProcessing(ctx context.Context)

	// Shutdown 処理を停止
	Shutdown()
}

// AudioBatcher 音声データのバッチ化インターフェース
type AudioBatcher interface {
	// AddAudioData 音声データをバッファに追加
	AddAudioData(clientID string, audioData []byte)

	// GetBatchReady 完成したバッチのチャネルを取得
	GetBatchReady() <-chan *model.AudioBatch

	// StartPeriodicFlush 定期フラッシュを開始
	StartPeriodicFlush(ctx context.Context)
}
