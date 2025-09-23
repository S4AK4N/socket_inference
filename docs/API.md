# 📡 API仕様書

Socket Inference システムの通信プロトコル仕様

## 🔌 WebSocket API

### エンドポイント
```
ws://localhost:8080/audio
```

### 接続時ヘッダー
```http
X-Client-ID: string  # クライアント識別ID（任意）
```

### メッセージフォーマット

#### 音声データ送信（クライアント → サーバー）
```
Type: Binary Message
Content: Raw audio data (bytes)
Max Size: 制限なし（推奨: 1KB - 8KB per chunk）
```

#### 接続例（JavaScript）
```javascript
const ws = new WebSocket('ws://localhost:8080/audio', [], {
    headers: {
        'X-Client-ID': 'client-001'
    }
});

// 音声データ送信
const audioChunk = new Uint8Array([/* audio data */]);
ws.send(audioChunk);
```

#### 接続例（Go）
```go
conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/audio", &websocket.DialOptions{
    HTTPHeader: map[string][]string{
        "X-Client-ID": {"client-001"},
    },
})

// 音声データ送信
audioData := []byte("audio chunk data")
err = conn.Write(ctx, websocket.MessageBinary, audioData)
```

### 接続ライフサイクル

1. **接続確立**
   - クライアントがWebSocket接続を確立
   - サーバーがクライアントを登録
   - ログ出力: `クライアント {ID} が接続しました`

2. **音声データストリーミング**
   - クライアントがバイナリメッセージで音声データを送信
   - サーバーがバッチ処理のためデータを蓄積
   - バッチ条件達成時に推論処理実行

3. **接続終了**
   - クライアントまたはサーバーが接続を閉じる
   - サーバーがクライアントの登録を解除
   - ログ出力: `クライアント {ID} が切断されました`

## 🔄 バッチ処理仕様

### バッチ生成条件
```
条件A: チャンク数 >= BATCH_SIZE（デフォルト: 10）
条件B: 最後のチャンク受信から FLUSH_TIMEOUT 経過（デフォルト: 2秒）
```

### バッチ形式
```go
type AudioBatch struct {
    ClientID  string    `json:"client_id"`
    AudioData [][]byte  `json:"audio_data"`
    Timestamp time.Time `json:"timestamp"`
    BatchSize int       `json:"batch_size"`
}
```

## 📨 gRPC API（Infrastructure Layer）

### サービス定義（プレースホルダー）
```protobuf
service InferenceService {
    rpc ProcessAudio(AudioRequest) returns (AudioResponse);
}

message AudioRequest {
    string client_id = 1;
    repeated bytes audio_chunks = 2;
    int64 timestamp = 3;
}

message AudioResponse {
    string client_id = 1;
    string result = 2;
    int32 status_code = 3;
    string message = 4;
}
```

### エンドポイント
```
デフォルト: localhost:50051
環境変数: GRPC_SERVER
```

## 🔧 HTTP管理API

### ヘルスチェック
```http
GET /audio
Response: WebSocket Upgrade または 400 Bad Request
```

### サーバー情報
```http
GET /
Response: 404 Not Found（WebSocketのみ対応）
```

## 🚨 エラーハンドリング

### WebSocket接続エラー
```
1001: 正常終了
1006: 異常終了（ネットワークエラー）
1000: 通常の切断
```

### 音声データエラー
- **送信タイムアウト**: 5秒でタイムアウト
- **データサイズ**: 制限なし（メモリ制約による）
- **フォーマット**: バイナリデータのみ受信

## 📊 監視・ログ

### 接続ログ
```
2024/XX/XX XX:XX:XX クライアント {client-id} が接続しました
2024/XX/XX XX:XX:XX クライアント {client-id} が切断されました
```

### バッチ処理ログ
```
2024/XX/XX XX:XX:XX バッチ準備完了: クライアント={client-id}, チャンク数={count}
2024/XX/XX XX:XX:XX 推論サーバーにバッチ送信: クライアント={client-id}, チャンク数={count}
```

### エラーログ
```
2024/XX/XX XX:XX:XX 音声データ読み取りエラー: {error}
2024/XX/XX XX:XX:XX 推論リクエスト失敗: {error}
```

## 🔒 セキュリティ考慮事項

### 本番環境での注意点
```go
// 開発用設定（本番では変更必要）
&websocket.AcceptOptions{
    InsecureSkipVerify: true, // 本番ではOriginチェック必須
}
```

### 推奨設定
- **Origin検証**: クロスオリジンリクエスト制御
- **接続数制限**: `MAX_CLIENTS` での同時接続制御
- **レート制限**: チャンク送信頻度の制限
- **データ検証**: 音声データの形式・サイズ検証

## 📈 パフォーマンス特性

### 推奨パラメータ
```bash
# 低レイテンシ重視
BATCH_SIZE=5 FLUSH_TIMEOUT=500ms

# スループット重視  
BATCH_SIZE=50 FLUSH_TIMEOUT=5s

# バランス型
BATCH_SIZE=10 FLUSH_TIMEOUT=2s
```

### 制限事項
- **最大同時接続**: 環境変数 `MAX_CLIENTS`（デフォルト: 100）
- **メモリ使用量**: バッチサイズ × クライアント数 × チャンクサイズ
- **ネットワーク帯域**: クライアント送信頻度に依存