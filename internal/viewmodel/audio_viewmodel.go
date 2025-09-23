package viewmodel

import (
	"context"
	"log"
	"time"

	"socket_inference/internal/infrastructure/interfaces"
	"socket_inference/internal/model"
	"socket_inference/internal/viewmodel/inference"
)

// AudioViewModel 音声処理と推論を管理
type AudioViewModel struct {
	clients map[*model.AudioClient]bool
	batcher *AudioBatcher
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewAudioViewModel 新しいAudioViewModelを作成
func NewAudioViewModel(inferenceClient interfaces.InferenceClient) *AudioViewModel {
	ctx, cancel := context.WithCancel(context.Background())
	batcher := NewAudioBatcher(10, 2*time.Second)

	vm := &AudioViewModel{
		clients: make(map[*model.AudioClient]bool),
		batcher: batcher,
		ctx:     ctx,
		cancel:  cancel,
	}

	// バッチャーにinferenceClientを渡すための処理を開始
	vm.startInferenceProcessing(inferenceClient)

	// バッチング処理を開始
	vm.startProcessing()

	return vm
}

// startInferenceProcessing 推論処理を開始
func (vm *AudioViewModel) startInferenceProcessing(inferenceClient interfaces.InferenceClient) {
	// 推論マネージャーを作成
	inferenceManager := vm.createInferenceManager(inferenceClient)

	// バッチ処理を開始
	go vm.batcher.StartBatching(vm.ctx)

	// 推論処理を開始
	inferenceManager.StartProcessing(vm.ctx, vm.batcher.GetBatchChannel())

	log.Println("推論処理を開始しました")
}

// createInferenceManager 推論マネージャーを作成
func (vm *AudioViewModel) createInferenceManager(inferenceClient interfaces.InferenceClient) *inference.Manager {
	return inference.NewManager(inferenceClient).(*inference.Manager)
}

// RegisterClient 新しい音声クライアントを登録
func (vm *AudioViewModel) RegisterClient(client *model.AudioClient) {
	vm.clients[client] = true
	log.Printf("音声クライアント接続: %s", client.ClientID)
}

// UnregisterClient 音声クライアントの登録を解除
func (vm *AudioViewModel) UnregisterClient(client *model.AudioClient) {
	if _, ok := vm.clients[client]; ok {
		delete(vm.clients, client)
		log.Printf("音声クライアント切断: %s", client.ClientID)
	}
}

// ProcessAudioData 受信した音声データを処理
func (vm *AudioViewModel) ProcessAudioData(clientID string, audioData []byte) {
	// 音声データをバッチャーに追加
	vm.batcher.AddAudioData(clientID, audioData)
}

// startProcessing バックグラウンド処理goroutineを開始
func (vm *AudioViewModel) startProcessing() {
	// バッチングの定期フラッシュを開始
	vm.batcher.StartPeriodicFlush(vm.ctx)

	// 推論処理用のgoroutineを開始
	go vm.processInferenceBatches()
}

// processInferenceBatches 完成したバッチの推論処理を実行
func (vm *AudioViewModel) processInferenceBatches() {
	for {
		select {
		case batch := <-vm.batcher.GetBatchReady():
			// 音声前処理を実行
			processedBatch := vm.preprocessAudioBatch(batch)

			// 推論サーバーに送信
			vm.sendToInferenceServer(processedBatch)

		case <-vm.ctx.Done():
			return
		}
	}
}

// preprocessAudioBatch 簡単な音声前処理を実行（プレースホルダー）
func (vm *AudioViewModel) preprocessAudioBatch(batch *model.AudioBatch) *model.AudioBatch {
	log.Printf("クライアント %s の音声バッチを前処理中: %d チャンク", batch.ClientID, batch.BatchSize)

	// 簡易的な前処理の雛形
	// 実際には音声フォーマット変換、正規化、フィルタリングなどを行う
	processedData := make([][]byte, len(batch.AudioData))
	for i, chunk := range batch.AudioData {
		// 例: 単純なデータコピー（実際の変換処理をここに追加）
		processedData[i] = make([]byte, len(chunk))
		copy(processedData[i], chunk)
	}

	return &model.AudioBatch{
		ClientID:  batch.ClientID,
		AudioData: processedData,
		Timestamp: batch.Timestamp,
		BatchSize: batch.BatchSize,
	}
}

// sendToInferenceServer 処理済みバッチをgRPC推論サーバーに送信（プレースホルダー）
func (vm *AudioViewModel) sendToInferenceServer(batch *model.AudioBatch) {
	log.Printf("推論サーバーにバッチ送信: クライアント=%s, チャンク数=%d, タイムスタンプ=%v",
		batch.ClientID, batch.BatchSize, batch.Timestamp)

	// TODO: gRPCクライアントの実装
	// 将来的にここでgRPCを使って推論サーバーに送信
	// 例:
	// req := &pb.InferenceRequest{
	//     ClientId: batch.ClientID,
	//     AudioData: batch.AudioData,
	//     Timestamp: batch.Timestamp.Unix(),
	// }
	// resp, err := grpcClient.ProcessAudio(ctx, req)
}

// Shutdown AudioViewModelを正常に停止
func (vm *AudioViewModel) Shutdown() {
	vm.cancel()
}
