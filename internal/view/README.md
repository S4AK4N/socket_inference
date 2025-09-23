# View Layer Architecture

View層は3つのレイヤーで構成されています：

## 📋 Layer Structure

```
internal/view/
├── interfaces/     # 抽象レイヤー（インターフェース定義）
├── handlers/       # 実装レイヤー（具体的な処理）
└── server/         # 設定レイヤー（サーバー構成）
```

## 🎯 Layer Responsibilities

### 1. Interfaces Layer (抽象レイヤー)
- **目的**: プロトコル非依存の共通インターフェース定義
- **ファイル**: `audio_stream_interface.go`
- **内容**:
  - `AudioStreamHandler` - 音声ストリーミング処理インターフェース
  - `AudioViewModelInterface` - ViewModelとの連携インターフェース

### 2. Handlers Layer (実装レイヤー)
- **目的**: プロトコル固有の具体的な実装
- **構造**:
  ```
  02_handlers/
  ├── websocket/
  │   └── audio_stream_handler.go    # WebSocket実装
  └── grpc/                          # 将来のgRPC実装
      └── audio_stream_handler.go
  ```
- **特徴**: 各プロトコルごとに独立した実装

### 3. Server Layer (設定レイヤー)
- **目的**: HTTPサーバーの設定とルーティング
- **ファイル**: `server.go`
- **責務**: サーバー起動、ルート設定、ミドルウェア管理

## 🔄 Dependency Flow

```
Interfaces (抽象)
    ↓ implements
Handlers (具象実装)
    ↓ uses
Server (設定・起動)
```

## 🚀 Extension Pattern

新しいプロトコル追加時:
1. `01_interfaces/` - インターフェースは変更不要
2. `02_handlers/new_protocol/` - 新しいプロトコル実装を追加
3. `03_server/` - 必要に応じてルーティング追加

この構造により、プロトコルに依存しない拡張可能なアーキテクチャを実現しています。