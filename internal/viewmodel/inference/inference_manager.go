package inference

import (
	"context"
	"log"
	"time"

	"socket_inference/internal/infrastructure/interfaces"
	"socket_inference/internal/model"
	vmInterfaces "socket_inference/internal/viewmodel/interfaces"
)

// Manager 推論処理管理の実装
type Manager struct {
	preprocessor    vmInterfaces.AudioPreprocessor
	inferenceClient interfaces.InferenceClient // Infrastructure依存を注入
	resultChannel   chan *model.InferenceResponse
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewManager 新しい推論マネージャーを作成
func NewManager(inferenceClient interfaces.InferenceClient) vmInterfaces.InferenceManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		preprocessor:    NewPreprocessor(),
		inferenceClient: inferenceClient,
		resultChannel:   make(chan *model.InferenceResponse, 100),
		ctx:             ctx,
		cancel:          cancel,
	}
}

// ProcessBatch バッチを推論処理
func (im *Manager) ProcessBatch(batch *model.AudioBatch) (*model.InferenceResponse, error) {
	// 前処理を実行
	processedBatch, err := im.preprocessor.PreprocessBatch(batch)
	if err != nil {
		return nil, err
	}

	// Infrastructure層のクライアントを使用して推論実行
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := im.inferenceClient.SendBatchInferenceRequest(ctx, processedBatch)
	if err != nil {
		log.Printf("推論リクエスト失敗: %v", err)
		return nil, err
	}

	return response, nil
}

// StartProcessing バックグラウンド推論処理を開始
func (im *Manager) StartProcessing(ctx context.Context, batchChan <-chan *model.AudioBatch) {
	go func() {
		for {
			select {
			case batch := <-batchChan:
				response, err := im.ProcessBatch(batch)
				if err != nil {
					log.Printf("推論処理エラー: %v", err)
					continue
				}

				select {
				case im.resultChannel <- response:
					log.Printf("推論結果送信完了: クライアント=%s", response.ClientID)
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()
	log.Println("推論処理マネージャーを開始しました")
}

// GetResultChannel 推論結果のチャネルを取得
func (im *Manager) GetResultChannel() <-chan *model.InferenceResponse {
	return im.resultChannel
}

// Shutdown 推論処理を停止
func (im *Manager) Shutdown() {
	im.cancel()
	close(im.resultChannel)
	log.Println("推論処理マネージャーを停止しました")
}
