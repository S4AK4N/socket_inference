package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"socket_inference/internal/model"
	interfaces "socket_inference/internal/view/interfaces"

	"github.com/coder/websocket"
)

// AudioStreamHandler WebSocketを使用した音声ストリーミングハンドラー
type AudioStreamHandler struct {
	viewModel interfaces.AudioViewModelInterface
}

// NewAudioStreamHandler 新しいAudioStreamHandlerを作成
func NewAudioStreamHandler(viewModel interfaces.AudioViewModelInterface) *AudioStreamHandler {
	return &AudioStreamHandler{
		viewModel: viewModel,
	}
}

// HandleWebSocket 音声ストリーミング用のWebSocket接続を処理
func (h *AudioStreamHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // 開発用。本番はOriginチェックを！
	})
	if err != nil {
		log.Println("accept:", err)
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = "unknown"
	}

	client := &model.AudioClient{
		Conn:     c,
		ClientID: clientID,
	}

	// クライアントをViewModelに登録
	h.viewModel.RegisterClient(client)

	// 読み取りループ - クライアントからの音声データを受信
	ctx := r.Context()
	defer func() {
		h.viewModel.UnregisterClient(client)
		_ = client.Conn.Close(websocket.StatusNormalClosure, "bye")
	}()

	for {
		readCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		_, audioData, err := client.Conn.Read(readCtx)
		cancel()
		if err != nil {
			log.Printf("クライアント %s の音声データ読み取りエラー: %v", client.ClientID, err)
			return
		}

		// 音声データをViewModelに送信
		h.viewModel.ProcessAudioData(client.ClientID, audioData)
	}
}

// HandleConnection AudioStreamHandlerインターフェースの実装
func (h *AudioStreamHandler) HandleConnection(connectionData interface{}) error {
	// この実装はHTTPハンドラーとして使用されるため、
	// 直接的なインターフェース実装は将来のgRPC統合で使用予定
	return nil
}
