# 第7章 リポジトリと永続化

『現場で使える Go 言語ドメイン駆動設計』第7章のサンプルコード。

## 動作環境

- Go 1.25.0+
- PostgreSQL 16+ (統合テスト時のみ、Testcontainers 経由)
- Docker (統合テスト時のみ)

## 構成

- `internal/booking/domain/` — ドメイン層 (集約 + ID 値オブジェクト + リポジトリ interface)
- `internal/booking/infrastructure/` — インフラ層 (in-memory リポジトリ + Persistence Model + マッピング)
- `db/schema.sql` — bookings テーブル定義

## 実行

```bash
go test ./...                       # ユニットテスト (in-memory)
go test -race ./...                 # race detector 付き
go test -tags=integration ./...     # 統合テスト (要 Docker、本リポでは未実装)
```

## ライセンス

MIT License
