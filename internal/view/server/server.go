package server

import (
	"log"
	"net/http"

	"socket_inference/internal/view/handlers/websocket"
)

// Server HTTPサーバーを表現
type Server struct {
	audioHandler *websocket.AudioStreamHandler
}

// NewServer 新しいHTTPサーバーを作成
func NewServer(audioHandler *websocket.AudioStreamHandler) *Server {
	return &Server{
		audioHandler: audioHandler,
	}
}

// SetupRoutes HTTPルートを設定
func (s *Server) SetupRoutes() {
	http.HandleFunc("/audio", s.audioHandler.HandleWebSocket)
}

// Start HTTPサーバーを開始
func (s *Server) Start(addr string) error {
	s.SetupRoutes()
	log.Printf("音声ストリーミングサーバーがリスニング中: %s", addr)
	log.Printf("接続先: ws://%s/audio", addr)
	return http.ListenAndServe(addr, nil)
}
