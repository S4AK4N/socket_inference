package coordinator

import (
	"context"
	"log"
	"time"

	"socket_inference/internal/infrastructure/interfaces"
	"socket_inference/internal/model"
	"socket_inference/internal/viewmodel/audio"
	"socket_inference/internal/viewmodel/client"
	"socket_inference/internal/viewmodel/inference"
	vmInterfaces "socket_inference/internal/viewmodel/interfaces"
)

// AudioViewModel 軽量化された全体調整ViewModelの実装
type AudioViewModel struct {
	clientManager    vmInterfaces.ClientManager
	audioProcessor   vmInterfaces.AudioProcessor
	inferenceManager vmInterfaces.InferenceManager
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewAudioViewModel 新しいAudioViewModelを作成
func NewAudioViewModel(inferenceClient interfaces.InferenceClient) *AudioViewModel {
	ctx, cancel := context.WithCancel(context.Background())

	// 各コンポーネントを初期化
	clientManager := client.NewManager()
	audioProcessor := audio.NewProcessor(10, 2*time.Second)
	inferenceManager := inference.NewManager(inferenceClient)

	vm := &AudioViewModel{
		clientManager:    clientManager,
		audioProcessor:   audioProcessor,
		inferenceManager: inferenceManager,
		ctx:              ctx,
		cancel:           cancel,
	}

	// バックグラウンド処理を開始
	vm.startProcessing()

	return vm
}

// RegisterClient 新しい音声クライアントを登録
func (vm *AudioViewModel) RegisterClient(client *model.AudioClient) {
	vm.clientManager.RegisterClient(client)
}

// UnregisterClient 音声クライアントの登録を解除
func (vm *AudioViewModel) UnregisterClient(client *model.AudioClient) {
	vm.clientManager.UnregisterClient(client)
}

// ProcessAudioData 受信した音声データを処理
func (vm *AudioViewModel) ProcessAudioData(clientID string, audioData []byte) {
	vm.audioProcessor.ProcessAudioData(clientID, audioData)
}

// startProcessing バックグラウンド処理を開始
func (vm *AudioViewModel) startProcessing() {
	// 音声処理を開始
	vm.audioProcessor.StartProcessing(vm.ctx)

	// 推論処理を開始
	vm.inferenceManager.StartProcessing(vm.ctx, vm.audioProcessor.GetBatchReady())

	// 推論結果の処理を開始
	go vm.processInferenceResults()

	log.Println("AudioViewModel: 全てのバックグラウンド処理を開始しました")
}

// processInferenceResults 推論結果の処理
func (vm *AudioViewModel) processInferenceResults() {
	for {
		select {
		case result := <-vm.inferenceManager.GetResultChannel():
			// TODO: 推論結果をクライアントに送信
			log.Printf("推論結果受信: クライアント=%s, 結果=%s, 信頼度=%.2f",
				result.ClientID, result.Result, result.Confidence)
		case <-vm.ctx.Done():
			return
		}
	}
}

// Shutdown AudioViewModelを正常に停止
func (vm *AudioViewModel) Shutdown() {
	log.Println("AudioViewModel: シャットダウンを開始します")

	vm.cancel()
	vm.audioProcessor.Shutdown()
	vm.inferenceManager.Shutdown()

	log.Println("AudioViewModel: シャットダウンが完了しました")
}
