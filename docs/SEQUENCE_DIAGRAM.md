# システムシーケンス図

このドキュメントでは、WebSocket接続プールシステムの動作フローをシーケンス図で説明します。

## 1. 接続プール初期化シーケンス

```mermaid
sequenceDiagram
    participant App as アプリケーション
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Factory as ConnectionFactory
    participant Server as WebSocketサーバー

    App->>Pool: NewConnectionPool(config)
    Pool->>Manager: NewDefaultPoolManager()
    Pool->>Factory: NewWebSocketConnectionFactory(url)
    
    Note over Pool,Manager: プール初期化
    Pool->>Pool: 接続プールを作成
    Pool->>Pool: 待機キューを初期化
    
    Pool-->>App: プール準備完了
```

## 2. 通常の接続取得・返却シーケンス

```mermaid
sequenceDiagram
    participant Client as クライアント
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Factory as ConnectionFactory
    participant Server as WebSocketサーバー

    Client->>Pool: GetConnection()
    
    alt 利用可能な接続がプールにある
        Pool->>Manager: GetConnection(poolId)
        Manager->>Manager: 接続を検索
        Manager-->>Pool: 既存接続を返却
        Pool-->>Client: 接続を提供
    else プールが空の場合
        Pool->>Factory: CreateConnection()
        Factory->>Server: WebSocket接続確立
        Server-->>Factory: 接続成功
        Factory-->>Pool: 新しい接続
        Pool->>Manager: AddConnection(poolId, conn)
        Pool-->>Client: 新しい接続を提供
    end

    Note over Client: クライアントが接続を使用

    Client->>Pool: ReleaseConnection(conn)
    Pool->>Manager: ReturnConnection(poolId, conn)
    Manager->>Manager: 接続をプールに戻す
    Pool-->>Client: 返却完了
```

## 3. 接続待機システムシーケンス

```mermaid
sequenceDiagram
    participant C1 as クライアント1
    participant C2 as クライアント2
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Queue as 待機キュー

    Note over Pool,Manager: プールが満杯状態

    C1->>Pool: GetConnection()
    Pool->>Manager: GetConnection(poolId)
    Manager-->>Pool: nil (利用可能な接続なし)
    
    Pool->>Queue: 待機キューに追加
    Note over Pool,Queue: C1は待機状態

    C2->>Pool: ReleaseConnection(conn)
    Pool->>Manager: ReturnConnection(poolId, conn)
    
    Pool->>Queue: 待機中クライアントを確認
    Queue-->>Pool: C1が待機中
    Pool->>Queue: C1に接続を通知
    Queue-->>C1: 接続利用可能
    
    C1->>Pool: 接続を取得
    Pool-->>C1: 接続を提供
```

## 4. 高負荷時の制御フローシーケンス

```mermaid
sequenceDiagram
    participant MC as 複数クライアント
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Limiter as 接続制限
    participant Monitor as 統計モニター

    MC->>Pool: 同時接続リクエスト(100個)
    
    loop 各リクエストに対して
        Pool->>Limiter: 接続制限チェック
        
        alt 制限内の場合
            Pool->>Manager: GetConnection()
            Manager-->>Pool: 接続またはnull
            
            alt 接続利用可能
                Pool-->>MC: 接続提供
                Pool->>Monitor: 接続統計更新
            else 接続待ち
                Pool->>Pool: 待機キューに追加
                Note over Pool: タイムアウト監視開始
            end
        else 制限超過
            Pool-->>MC: 接続拒否
            Pool->>Monitor: 拒否統計更新
        end
    end

    Note over Monitor: 統計情報の収集・分析
```

## 5. エラーハンドリング・クリーンアップシーケンス

```mermaid
sequenceDiagram
    participant Client as クライアント
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Conn as WebSocket接続
    participant Cleanup as クリーンアップ

    Client->>Pool: GetConnection()
    Pool->>Manager: GetConnection()
    Manager-->>Pool: 接続
    Pool-->>Client: 接続提供

    Note over Client,Conn: 通信エラー発生
    Conn->>Pool: 接続エラー通知
    
    Pool->>Manager: RemoveConnection(poolId, conn)
    Manager->>Manager: 不正な接続を削除
    
    Pool->>Cleanup: CleanupConnection(conn)
    Cleanup->>Conn: Close()
    Cleanup-->>Pool: クリーンアップ完了

    alt クライアントが再試行
        Client->>Pool: GetConnection() (再試行)
        Pool->>Pool: 新しい接続を作成
        Pool-->>Client: 新しい接続を提供
    end
```

## 6. システム終了時のクリーンアップシーケンス

```mermaid
sequenceDiagram
    participant App as アプリケーション
    participant Pool as ConnectionPool
    participant Manager as PoolManager
    participant Conns as 全接続
    participant Queue as 待機キュー

    App->>Pool: Close()
    
    Pool->>Queue: 待機中のクライアントに通知
    Queue-->>Pool: 全ての待機を解除
    
    Pool->>Manager: GetAllConnections()
    Manager-->>Pool: 全接続リスト
    
    loop 各接続に対して
        Pool->>Conns: Close()
        Conns-->>Pool: 接続終了
    end
    
    Pool->>Manager: Cleanup()
    Manager-->>Pool: マネージャー終了
    
    Pool->>Pool: リソース解放
    Pool-->>App: システム終了完了
```

## アーキテクチャー概要

このシステムは以下の主要コンポーネントで構成されています：

- **ConnectionPool**: メインオーケストレーター、待機システム管理
- **PoolManager**: 実際の接続プール管理、スレッドセーフ
- **ConnectionFactory**: WebSocket接続の作成を抽象化
- **統計システム**: パフォーマンス監視と分析
- **クリーンアップシステム**: リソース管理と終了処理

各コンポーネントは疎結合で設計されており、Clean Architectureの原則に従って実装されています。
