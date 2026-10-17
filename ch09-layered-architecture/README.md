# 第9章 レイヤー化アーキテクチャとパッケージ分割

『現場で使える Go 言語ドメイン駆動設計』第9章のサンプルコード。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください）

## 実行

```bash
go test ./...                       # ユニットテスト
go test -race ./...                 # race detector
go test -cover ./...                # カバレッジ
go build ./cmd/api                  # API サーバビルド
go run ./cmd/api                    # API サーバ起動 (:8080)
```

## 構成

```text
ch09-layered-architecture/
  cmd/api/main.go              # 全レイヤーのワイヤリング + graceful shutdown
  internal/
    order/
      domain/                  # ドメイン層 (エンティティ・値オブジェクト・Repository interface)
      application/             # アプリケーション層 (ユースケース・クエリ・ApplicationError)
      infrastructure/          # インフラ層 (in-memory リポジトリ)
      presentation/            # プレゼンテーション層 (HTTP ハンドラ・ミドルウェア・DTO)
```

依存方向は内向き一方向: Presentation → Application → Domain ← Infrastructure。

## API エンドポイント

- `POST /orders` — 注文作成
- `GET /orders/{id}` — 注文詳細取得
- `GET /orders?customer_id=cus_001` — 顧客別注文一覧

## トピック

- `internal/` パッケージでアクセス制御
- Higher-Order Handler（ハンドラ生成関数で依存注入）
- Go 1.22 の `net/http` 強化ルーティング
- ミドルウェアの合成（RequestID / Logging / Recovery）
- `httptest` でのハンドラ + ミドルウェアテスト
- Spring の `@RestController` + DI コンテナとの対比

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
