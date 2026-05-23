# 第7章 リポジトリと永続化

『現場で使える Go 言語ドメイン駆動設計』第7章のサンプルコード。

## 動作環境

- Go 1.25.0 以上
- PostgreSQL 16 以上（統合テストを動かす場合のみ、Testcontainers 経由）
- Docker（統合テストを動かす場合のみ）

## 実行

```bash
go test ./...                       # ユニットテスト (in-memory リポジトリ)
go test -race ./...                 # race detector 付き
go build -tags=pg ./...             # PostgreSQL 実装骨子のビルド確認
```

PostgreSQL / Testcontainers を使った統合テストは本文 7.5-7.8 のコードを写経して読者が組み立てる構成です。本サンプル repo には骨子のみ収録しています。

## 構成

- `internal/booking/domain/` — ドメイン層（集約 + ID 値オブジェクト + Repository interface）
- `internal/booking/infrastructure/` — インフラ層（in-memory リポジトリ）
- `internal/booking/infrastructure/booking_row.go`, `booking_mapper.go` — PostgreSQL リポジトリ用の Persistence Model + マッピング（`//go:build pg` で隔離、骨子のみ）
- `db/schema.sql` — bookings テーブル定義（本文の sqlc 用）

## トピック

- consumer-side interface（ドメイン層に Repository interface を置く）
- ドメインモデル ↔ Persistence Model の分離
- in-memory リポジトリでテストダブルを作る（モックではなくフェイク）
- 楽観的ロック（version スナップショット）の実装
- `pgx.ErrNoRows` → `domain.ErrNotFound` への変換
- `context.Context` にトランザクションを乗せる方式
- Testcontainers でのリアル DB 統合テスト（本文参照）

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
