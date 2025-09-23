package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"socket_inference/internal/infrastructure/grpc"
	"socket_inference/internal/view/handlers/websocket"
	"socket_inference/internal/view/server"
	"socket_inference/internal/viewmodel/coordinator"
)

func main() {
	// Infrastructure層の実装を作成
	grpcClient := grpc.NewInferenceClient("localhost:50051", 30*time.Second)

	// ViewModelを作成（Infrastructure実装を注入）
	audioViewModel := coordinator.NewAudioViewModel(grpcClient)
	defer audioViewModel.Shutdown()

	// Viewを作成
	audioHandler := websocket.NewAudioStreamHandler(audioViewModel)
	httpServer := server.NewServer(audioHandler)

	// 正常なシャットダウンのためのシグナルハンドリング
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// サーバーを別のgoroutineで開始
	go func() {
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatalf("サーバー起動失敗: %v", err)
		}
	}()

	// シャットダウンシグナルを待機
	<-stop
	log.Println("サーバーを停止中...")
}
