package audio

import (
	"context"
	"log"
	"sync"
	"time"

	"socket_inference/internal/model"
)

// AudioBatcher 推論処理用の音声データバッチ化を処理
type AudioBatcher struct {
	mu           sync.Mutex
	audioBuffer  map[string][][]byte    // clientID -> audio chunks
	batchSize    int                    // バッチサイズ
	flushTimeout time.Duration          // フラッシュタイムアウト
	batchReady   chan *model.AudioBatch // 完成したバッチを送信するチャネル
	lastFlush    map[string]time.Time   // clientID -> 最後のフラッシュ時間
}

// NewAudioBatcher 新しいAudioBatcherを作成
func NewAudioBatcher(batchSize int, flushTimeout time.Duration) *AudioBatcher {
	return &AudioBatcher{
		audioBuffer:  make(map[string][][]byte),
		batchSize:    batchSize,
		flushTimeout: flushTimeout,
		batchReady:   make(chan *model.AudioBatch, 100),
		lastFlush:    make(map[string]time.Time),
	}
}

// AddAudioData 音声データをバッファに追加し、バッチ準備状況をチェック
func (ab *AudioBatcher) AddAudioData(clientID string, audioData []byte) {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	// 音声データをバッファに追加
	ab.audioBuffer[clientID] = append(ab.audioBuffer[clientID], audioData)

	// バッチサイズに達したかチェック
	if len(ab.audioBuffer[clientID]) >= ab.batchSize {
		ab.flushBatch(clientID)
	}
}

// flushBatch バッチを作成し、準備完了チャネルに送信
func (ab *AudioBatcher) flushBatch(clientID string) {
	if len(ab.audioBuffer[clientID]) == 0 {
		return
	}

	batch := &model.AudioBatch{
		ClientID:  clientID,
		AudioData: make([][]byte, len(ab.audioBuffer[clientID])),
		Timestamp: time.Now(),
		BatchSize: len(ab.audioBuffer[clientID]),
	}

	// データをコピー
	copy(batch.AudioData, ab.audioBuffer[clientID])

	// バッファをクリア
	ab.audioBuffer[clientID] = nil
	ab.lastFlush[clientID] = time.Now()

	// バッチを送信チャネルに送る
	select {
	case ab.batchReady <- batch:
		log.Printf("クライアント %s のバッチ準備完了: %d 音声チャンク", clientID, batch.BatchSize)
	default:
		log.Printf("バッチチャネルが満杯、クライアント %s のバッチを破棄", clientID)
	}
}

// StartPeriodicFlush 古いデータを定期的にフラッシュするgoroutineを開始
func (ab *AudioBatcher) StartPeriodicFlush(ctx context.Context) {
	ticker := time.NewTicker(ab.flushTimeout)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				ab.flushOldBatches()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// flushOldBatches 長時間待機しているバッチをフラッシュ
func (ab *AudioBatcher) flushOldBatches() {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	now := time.Now()
	for clientID, lastFlush := range ab.lastFlush {
		if now.Sub(lastFlush) > ab.flushTimeout && len(ab.audioBuffer[clientID]) > 0 {
			log.Printf("クライアント %s の古いバッチをフラッシュ（タイムアウト）", clientID)
			ab.flushBatch(clientID)
		}
	}
}

// GetBatchReady バッチ準備完了チャネルを返す
func (ab *AudioBatcher) GetBatchReady() <-chan *model.AudioBatch {
	return ab.batchReady
}
