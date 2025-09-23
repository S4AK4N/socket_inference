package inference

import (
	"log"

	"socket_inference/internal/model"
	"socket_inference/internal/viewmodel/interfaces"
)

// Preprocessor 音声前処理の実装
type Preprocessor struct {
	params map[string]interface{}
}

// NewPreprocessor 新しい前処理器を作成
func NewPreprocessor() interfaces.AudioPreprocessor {
	return &Preprocessor{
		params: make(map[string]interface{}),
	}
}

// PreprocessBatch 音声バッチの前処理
func (ap *Preprocessor) PreprocessBatch(batch *model.AudioBatch) (*model.AudioBatch, error) {
	log.Printf("クライアント %s の音声バッチを前処理中: %d チャンク", batch.ClientID, batch.BatchSize)

	// 前処理済みデータの作成
	processedData := make([][]byte, len(batch.AudioData))
	for i, chunk := range batch.AudioData {
		processedData[i] = make([]byte, len(chunk))
		copy(processedData[i], chunk)

		// TODO: 実際の前処理をここに追加
		// - 音声フォーマット変換
		// - 正規化
		// - ノイズ除去
		// - 特徴量抽出
	}

	return &model.AudioBatch{
		ClientID:  batch.ClientID,
		AudioData: processedData,
		Timestamp: batch.Timestamp,
		BatchSize: batch.BatchSize,
	}, nil
}

// SetPreprocessingParameters 前処理パラメータを設定
func (ap *Preprocessor) SetPreprocessingParameters(params map[string]interface{}) {
	ap.params = params
	log.Printf("前処理パラメータを更新しました: %+v", params)
}
