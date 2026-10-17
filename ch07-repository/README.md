# 第7章 リポジトリと永続化

『現場で使えるGo言語ドメイン駆動設計』第7章のサンプルコード。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください）
- sqlc v1.30.0（`db/*.sql` を変えて生成し直す場合のみ。`go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0`）
- PostgreSQL 16 以上（統合テストを動かす場合のみ、Testcontainers 経由）
- Docker（統合テストを動かす場合のみ）

## 実行

```bash
go test ./...                       # ユニットテスト (in-memory リポジトリ)
go test -race ./...                 # race detector 付き
go test -tags=integration ./...     # 統合テスト (Testcontainers で PostgreSQL を起動。要 Docker)
sqlc generate                       # db/*.sql を変えたら sqlcgen/ を作り直す (生成物はコミット済み)
```

## 構成

- `internal/booking/domain/` — ドメイン層（集約 + ID 値オブジェクト + Repository interface）
- `internal/booking/infrastructure/` — インフラ層（in-memory リポジトリ、PostgreSQL リポジトリ、トランザクションの伝播 `tx.go`）
- `internal/booking/infrastructure/booking_row.go`, `booking_mapper.go` — Persistence Model とマッピング
- `internal/booking/infrastructure/sqlcgen/` — `sqlc generate` の生成物
- `sqlc.yaml`, `db/schema.sql`, `db/queries.sql`, `db/list_queries.sql` — sqlc の設定・テーブル定義・クエリ
- `db/embed.go` — 統合テストの `MigrateSchema` が読む schema.sql を埋め込む

## トピック

- consumer-side interface（ドメイン層に Repository interface を置く）
- ドメインモデル ↔ Persistence Model の分離
- in-memory リポジトリでテストダブルを作る（モックではなくフェイク）
- 楽観的ロック（`UPDATE ... WHERE version = @expected_version` と in-memory での再現）
- `pgx.ErrNoRows` → `domain.ErrNotFound` への変換
- `context.Context` にトランザクションを乗せる方式
- Testcontainers でのリアル DB 統合テスト

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
