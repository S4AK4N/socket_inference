package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"socket_inference/pkg/connection_pool"
	"socket_inference/pkg/connection_pool/core"

	"github.com/coder/websocket"
)

// Config パラメータチューニング用設定
type Config struct {
	ServerURL       string        // サーバーURL
	ClientCount     int           // 並行クライアント数
	ChunksPerClient int           // クライアント毎の送信チャンク数
	ChunkInterval   time.Duration // チャンク送信間隔
	ChunkSize       int           // チャンクサイズ（バイト）
	TestDuration    time.Duration // テスト継続時間
	// 接続プール設定
	UseConnectionPool bool          // 接続プールを使用するか
	PoolSize          int           // 接続プールサイズ
	ConnectTimeout    time.Duration // 接続タイムアウト
	IdleTimeout       time.Duration // アイドルタイムアウト
}

// LoadConfig 環境変数から設定を読み込み
func LoadConfig() *Config {
	config := &Config{
		ServerURL:         getEnv("SERVER_URL", "ws://localhost:8080/audio"),
		ClientCount:       getEnvInt("CLIENT_COUNT", 3),
		ChunksPerClient:   getEnvInt("CHUNKS_PER_CLIENT", 15),
		ChunkInterval:     getEnvDuration("CHUNK_INTERVAL", "100ms"),
		ChunkSize:         getEnvInt("CHUNK_SIZE", 1024),
		TestDuration:      getEnvDuration("TEST_DURATION", "10s"),
		UseConnectionPool: getEnvBool("USE_CONNECTION_POOL", true),
		PoolSize:          getEnvInt("POOL_SIZE", 50),
		ConnectTimeout:    getEnvDuration("CONNECT_TIMEOUT", "10s"),
		IdleTimeout:       getEnvDuration("IDLE_TIMEOUT", "5m"),
	}
	return config
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

// getEnvDuration 環境変数から時間取得
func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

// getEnvBool 環境変数からbool値取得
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// AudioClient 音声ストリーミングクライアント
type AudioClient struct {
	ID             int
	config         *Config
	conn           *websocket.Conn
	pooledConn     *core.PooledConnection
	connectionPool *connection_pool.ConnectionPool
	stats          *ClientStats
}

// ClientStats クライアント統計情報
type ClientStats struct {
	SentChunks int
	BytesSent  int64
	ErrorCount int
	Duration   time.Duration
	StartTime  time.Time
}

// NewAudioClient 新しい音声クライアントを作成
func NewAudioClient(id int, config *Config, pool *connection_pool.ConnectionPool) *AudioClient {
	return &AudioClient{
		ID:             id,
		config:         config,
		connectionPool: pool,
		stats: &ClientStats{
			StartTime: time.Now(),
		},
	}
}

// Connect WebSocketサーバーに接続
func (c *AudioClient) Connect(ctx context.Context) error {
	if c.config.UseConnectionPool && c.connectionPool != nil {
		// 接続プールから接続を取得
		pooledConn, err := c.connectionPool.Get(ctx)
		if err != nil {
			return fmt.Errorf("接続プールから接続取得失敗: %v", err)
		}
		c.pooledConn = pooledConn
		c.conn = pooledConn.Conn

		return nil
	} else {
		// 従来の直接接続
		u, err := url.Parse(c.config.ServerURL)
		if err != nil {
			return fmt.Errorf("URL解析エラー: %v", err)
		}

		conn, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
			HTTPHeader: map[string][]string{
				"X-Client-ID": {fmt.Sprintf("tune-client-%03d", c.ID)},
			},
		})
		if err != nil {
			return fmt.Errorf("WebSocket接続エラー: %v", err)
		}

		c.conn = conn
		return nil
	}
}

// SendAudioData 音声データストリーミング実行
func (c *AudioClient) SendAudioData(ctx context.Context) error {
	defer func() {
		if c.config.UseConnectionPool && c.pooledConn != nil {
			// 接続プールに接続を返却
			c.connectionPool.Put(c.pooledConn)
		} else {
			// 通常の接続は閉じる
			c.conn.Close(websocket.StatusNormalClosure, "テスト完了")
		}
	}()

	// テスト期間のタイムアウト設定
	testCtx, cancel := context.WithTimeout(ctx, c.config.TestDuration)
	defer cancel()

	chunkData := make([]byte, c.config.ChunkSize)
	// チャンクデータを初期化（クライアントIDベース）
	for i := range chunkData {
		chunkData[i] = byte((c.ID + i) % 256)
	}

	ticker := time.NewTicker(c.config.ChunkInterval)
	defer ticker.Stop()

	chunkCount := 0
	for {
		select {
		case <-testCtx.Done():
			// テスト期間終了
			c.stats.Duration = time.Since(c.stats.StartTime)
			return nil
		case <-ticker.C:
			if chunkCount >= c.config.ChunksPerClient {
				// 指定チャンク数送信完了
				c.stats.Duration = time.Since(c.stats.StartTime)
				return nil
			}

			// チャンクにシーケンス番号を埋め込み
			copy(chunkData[:4], []byte(fmt.Sprintf("%04d", chunkCount)))

			writeCtx, writeCancel := context.WithTimeout(testCtx, 5*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageBinary, chunkData)
			writeCancel()

			if err != nil {
				c.stats.ErrorCount++
				log.Printf("❌ [Client-%03d] チャンク %d 送信失敗: %v", c.ID, chunkCount, err)
				continue
			}

			c.stats.SentChunks++
			c.stats.BytesSent += int64(len(chunkData))
			chunkCount++

			if chunkCount%10 == 0 {
				log.Printf("📤 [Client-%03d] %d チャンク送信完了", c.ID, chunkCount)
			}
		}
	}
}

// PrintStats 統計情報を出力
func (c *AudioClient) PrintStats() {
	throughput := float64(c.stats.BytesSent) / c.stats.Duration.Seconds() / 1024 // KB/s
	fmt.Printf("📊 [Client-%03d] 統計: チャンク=%d, バイト=%d, エラー=%d, 期間=%.2fs, スループット=%.2fKB/s\n",
		c.ID, c.stats.SentChunks, c.stats.BytesSent, c.stats.ErrorCount,
		c.stats.Duration.Seconds(), throughput)
}

func main() {
	fmt.Println("🎯 パラメータチューニング用音声ストリーミングテスト")
	fmt.Println("================================================")

	// 設定読み込み
	config := LoadConfig()

	fmt.Printf("⚙️  設定:\n")
	fmt.Printf("   - サーバーURL: %s\n", config.ServerURL)
	fmt.Printf("   - 並行クライアント数: %d\n", config.ClientCount)
	fmt.Printf("   - クライアント毎チャンク数: %d\n", config.ChunksPerClient)
	fmt.Printf("   - チャンク送信間隔: %v\n", config.ChunkInterval)
	fmt.Printf("   - チャンクサイズ: %d bytes\n", config.ChunkSize)
	fmt.Printf("   - テスト継続時間: %v\n", config.TestDuration)
	fmt.Printf("   - 接続プール使用: %t\n", config.UseConnectionPool)
	if config.UseConnectionPool {
		fmt.Printf("   - プールサイズ: %d\n", config.PoolSize)
		fmt.Printf("   - 接続タイムアウト: %v\n", config.ConnectTimeout)
		fmt.Printf("   - アイドルタイムアウト: %v\n", config.IdleTimeout)
	}
	fmt.Println()

	// 接続プール初期化（必要な場合）
	var pool *connection_pool.ConnectionPool
	if config.UseConnectionPool {
		poolConfig := core.PoolConfig{
			MaxPoolSize:     config.PoolSize,
			ConnectTimeout:  config.ConnectTimeout,
			IdleTimeout:     config.IdleTimeout,
			MaxLifetime:     30 * time.Minute,
			CleanupInterval: 1 * time.Minute,
			ServerURL:       config.ServerURL,
		}
		pool = connection_pool.NewConnectionPool(poolConfig)
		defer pool.Shutdown()

		fmt.Printf("🔗 接続プール初期化完了 (最大接続数: %d)\n", config.PoolSize)
	}

	// 全体のテストコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), config.TestDuration+10*time.Second)
	defer cancel()

	// 並行クライアント実行
	var wg sync.WaitGroup
	clients := make([]*AudioClient, config.ClientCount)

	startTime := time.Now()
	fmt.Printf("🚀 %d 個のクライアントでテスト開始...\n", config.ClientCount)

	for i := 0; i < config.ClientCount; i++ {
		wg.Add(1)
		clients[i] = NewAudioClient(i+1, config, pool)

		go func(client *AudioClient) {
			defer wg.Done()

			// 接続
			if err := client.Connect(ctx); err != nil {
				log.Printf("❌ [Client-%03d] 接続失敗: %v", client.ID, err)
				return
			}
			fmt.Printf("✅ [Client-%03d] 接続成功\n", client.ID)

			// 音声データ送信
			if err := client.SendAudioData(ctx); err != nil {
				log.Printf("❌ [Client-%03d] データ送信エラー: %v", client.ID, err)
				return
			}
		}(clients[i])

		// 接続間隔を少しずらす
		time.Sleep(50 * time.Millisecond)
	}

	// 全クライアント完了待ち
	wg.Wait()
	totalDuration := time.Since(startTime)

	fmt.Println("\n📋 テスト結果サマリー:")
	fmt.Println("========================")

	// 個別統計
	totalChunks := 0
	totalBytes := int64(0)
	totalErrors := 0

	for _, client := range clients {
		client.PrintStats()
		totalChunks += client.stats.SentChunks
		totalBytes += client.stats.BytesSent
		totalErrors += client.stats.ErrorCount
	}

	// 全体統計
	overallThroughput := float64(totalBytes) / totalDuration.Seconds() / 1024 // KB/s
	fmt.Printf("\n🏆 全体統計:\n")
	fmt.Printf("   - 総チャンク数: %d\n", totalChunks)
	fmt.Printf("   - 総バイト数: %d (%.2f MB)\n", totalBytes, float64(totalBytes)/1024/1024)
	fmt.Printf("   - 総エラー数: %d\n", totalErrors)
	fmt.Printf("   - 実行時間: %.2f秒\n", totalDuration.Seconds())
	fmt.Printf("   - 全体スループット: %.2f KB/s\n", overallThroughput)
	fmt.Printf("   - 平均クライアントスループット: %.2f KB/s\n", overallThroughput/float64(config.ClientCount))

	fmt.Println("\n✨ パラメータチューニングテスト完了")
}
