package audio

import (
	"context"
	"log"
	"time"

	"socket_inference/internal/model"
	"socket_inference/internal/viewmodel/interfaces"
)

// Processor 音声データ処理の実装
type Processor struct {
	batcher interfaces.AudioBatcher
}

// NewProcessor 新しい音声プロセッサーを作成
func NewProcessor(batchSize int, flushTimeout time.Duration) interfaces.AudioProcessor {
	batcher := NewAudioBatcher(batchSize, flushTimeout)
	return &Processor{
		batcher: batcher,
	}
}

// ProcessAudioData 音声データを処理してバッチ化
func (p *Processor) ProcessAudioData(clientID string, audioData []byte) {
	p.batcher.AddAudioData(clientID, audioData)
}

// GetBatchReady 完成したバッチを受信するチャネルを取得
func (p *Processor) GetBatchReady() <-chan *model.AudioBatch {
	return p.batcher.GetBatchReady()
}

// StartProcessing バックグラウンド処理を開始
func (p *Processor) StartProcessing(ctx context.Context) {
	p.batcher.StartPeriodicFlush(ctx)
	log.Println("音声処理プロセッサーを開始しました")
}

// Shutdown 処理を停止
func (p *Processor) Shutdown() {
	log.Println("音声処理プロセッサーを停止します")
}
