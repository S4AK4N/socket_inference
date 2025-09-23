package model

import "time"

// InferenceRequest 推論サーバーへのリクエストを表現
// 音声データを推論処理するためのリクエストドメインモデル
type InferenceRequest struct {
	ClientID  string    `json:"client_id"`  // クライアント識別ID
	AudioData [][]byte  `json:"audio_data"` // 音声データ配列
	Timestamp time.Time `json:"timestamp"`  // リクエスト生成時刻
	BatchSize int       `json:"batch_size"` // バッチサイズ
}

// InferenceResponse 推論サーバーからのレスポンスを表現
// 推論処理結果を格納するレスポンスドメインモデル
type InferenceResponse struct {
	ClientID       string        `json:"client_id"`       // クライアント識別ID
	Result         string        `json:"result"`          // 推論結果
	Confidence     float64       `json:"confidence"`      // 推論の信頼度
	ProcessingTime time.Duration `json:"processing_time"` // 処理時間
}
