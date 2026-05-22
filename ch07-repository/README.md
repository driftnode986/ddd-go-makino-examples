# 第7章 リポジトリと永続化

『現場で使える Go 言語ドメイン駆動設計』第7章のサンプルコード。

## 動作環境

- Go 1.25.0+
- PostgreSQL 16+ (統合テスト時のみ、Testcontainers 経由)
- Docker (統合テスト時のみ)

## 構成

- `internal/booking/domain/` — ドメイン層 (集約 + ID 値オブジェクト + リポジトリ interface)
- `internal/booking/infrastructure/` — インフラ層 (in-memory リポジトリ)
- `internal/booking/infrastructure/booking_row.go`, `booking_mapper.go` — PostgreSQL リポジトリ用の Persistence Model + マッピング (`//go:build pg` で隔離、骨子のみ)
- `db/schema.sql` — bookings テーブル定義 (本文の sqlc 用)

## 実行

```bash
go test ./...                       # ユニットテスト (in-memory)
go test -race ./...                 # race detector 付き
go build -tags=pg ./...             # PG 実装骨子のビルド確認
```

PG / Testcontainers 統合テストは本文 7.5-7.8 を参照して読者が組み立てる構成です。

## ライセンス

MIT License
