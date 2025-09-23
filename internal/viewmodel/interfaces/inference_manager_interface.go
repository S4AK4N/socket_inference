package interfaces

import (
	"context"
	"socket_inference/internal/model"
)

// InferenceManager 推論処理管理のインターフェース
type InferenceManager interface {
	// ProcessBatch バッチを推論処理
	ProcessBatch(batch *model.AudioBatch) (*model.InferenceResponse, error)

	// StartProcessing バックグラウンド推論処理を開始
	StartProcessing(ctx context.Context, batchChan <-chan *model.AudioBatch)

	// GetResultChannel 推論結果のチャネルを取得
	GetResultChannel() <-chan *model.InferenceResponse

	// Shutdown 推論処理を停止
	Shutdown()
}

// AudioPreprocessor 音声前処理のインターフェース
type AudioPreprocessor interface {
	// PreprocessBatch 音声バッチの前処理
	PreprocessBatch(batch *model.AudioBatch) (*model.AudioBatch, error)

	// SetPreprocessingParameters 前処理パラメータを設定
	SetPreprocessingParameters(params map[string]interface{})
}
