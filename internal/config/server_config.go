package config

import (
	"os"
	"strconv"
	"time"
)

// ServerConfig サーバー設定
type ServerConfig struct {
	Port         string        // サーバーポート
	BatchSize    int           // 音声バッチサイズ
	FlushTimeout time.Duration // バッチフラッシュタイムアウト
	MaxClients   int           // 最大同時接続クライアント数
	BufferSize   int           // チャネルバッファサイズ
	GRPCServer   string        // gRPCサーバーアドレス
	GRPCTimeout  time.Duration // gRPCタイムアウト
}

// LoadServerConfig 環境変数からサーバー設定を読み込み
func LoadServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		BatchSize:    getEnvInt("BATCH_SIZE", 10),
		FlushTimeout: getEnvDuration("FLUSH_TIMEOUT", "2s"),
		MaxClients:   getEnvInt("MAX_CLIENTS", 100),
		BufferSize:   getEnvInt("BUFFER_SIZE", 100),
		GRPCServer:   getEnv("GRPC_SERVER", "localhost:50051"),
		GRPCTimeout:  getEnvDuration("GRPC_TIMEOUT", "30s"),
	}
}

// getEnv 環境変数取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 環境変数から整数取得
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration 環境変数から期間取得
func getEnvDuration(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// デフォルト値をパース
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
