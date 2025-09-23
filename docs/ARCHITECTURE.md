# Internal Architecture

Clean Architectureの依存関係方向に基づいた設計です。

## 🏗️ Layer Structure (抽象 → 具象)

```
internal/
├── model/              # Domain Layer (最抽象層)
├── viewmodel/          # Use Case Layer (ビジネスロジック)
├── view/               # Interface Adapter Layer (プレゼンテーション)
└── infrastructure/     # Infrastructure Layer (最具象層)
```

## 📋 Dependency Direction

```
Domain (model) ← Use Case (viewmodel) ← Interface Adapter (view) ← Infrastructure
```

**依存関係の原則:**
- 外側の層は内側の層に依存する
- 内側の層は外側の層を知らない
- 依存関係は常に内向き（抽象に向かう）

## 🎯 Layer Responsibilities

### 1. Domain Layer (model/)
- **責務**: ビジネスエンティティとルール
- **依存**: 他の層に依存しない（最抽象）
- **内容**: `AudioClient`, `AudioBatch`, `InferenceRequest/Response`

### 2. Use Case Layer (viewmodel/)
- **責務**: アプリケーション固有のビジネスロジック
- **依存**: Domain層のみに依存
- **内容**: クライアント管理、音声処理、推論管理、全体調整

### 3. Interface Adapter Layer (view/)
- **責務**: 外部との入出力データ変換
- **依存**: Domain層とUse Case層に依存
- **内容**: WebSocketハンドラー、HTTPサーバー

### 4. Infrastructure Layer (infrastructure/)
- **責務**: 外部システムとの実際の通信
- **依存**: 全ての層に依存可能（最具象）
- **内容**: gRPCクライアント、データベース、外部API

## 🔄 Interface Pattern

各層間は**インターフェースを通じて通信**し、具体的な実装に依存しません：

```go
// Use Case → Infrastructure (依存注入)
type InferenceManager interface { ... }
type InferenceClient interface { ... }

// Infrastructure実装をUse Caseに注入
manager := inference.NewManager(grpcClient)
```

## 🚀 Extension Strategy

新機能追加時:
1. **Domain**: 新しいエンティティを追加
2. **Use Case**: ビジネスロジックを実装
3. **Interface Adapter**: 入出力処理を追加
4. **Infrastructure**: 外部システム連携を実装

この順序により、Clean Architectureの原則を維持しながら拡張できます。