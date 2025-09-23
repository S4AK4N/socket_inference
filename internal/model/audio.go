package model

import "time"

// AudioBatch 推論処理用の音声データバッチを表現
// 音声データのバッチ化と管理を担当するドメインモデル
type AudioBatch struct {
	ClientID  string    `json:"client_id"`  // クライアント識別ID
	AudioData [][]byte  `json:"audio_data"` // 音声データ配列
	Timestamp time.Time `json:"timestamp"`  // バッチ生成時刻
	BatchSize int       `json:"batch_size"` // バッチサイズ
}
