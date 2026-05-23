# 第11章: 実プロジェクト構成例とテスト戦略

『現場で使える Go 言語ドメイン駆動設計』第11章のサンプルコード。

トレーニング予約 SaaS を 3 境界づけられたコンテキスト (Booking / Catalog / Identity) で構成し、`go.work` で 1 リポジトリにまとめる capstone 章のサンプルです。

## 動作環境

- Go 1.25.0+
- Docker (オプション、PostgreSQL を起動する場合のみ)

## ディレクトリ構成

```
ch11-practical-project/
  go.work                       # 4 module をまとめる
  docker-compose.yaml           # PostgreSQL 17
  Makefile
  shared/                       # 公表された言語 DTO
    dto/lesson_view.go
  booking/                      # Booking BC (フル実装サンプル)
    cmd/api/main.go
    internal/booking/
      domain/                   # Booking 集約 + repository interface
      application/              # ユースケース + クエリ
      infrastructure/memory/    # in-memory リポジトリ
      presentation/http/        # chi ハンドラ + httptest
      testdata/                 # Test Data Builder
  catalog/                      # Catalog BC (Lesson 集約のみ)
  identity/                     # Identity BC (Member 集約のみ)
```

## 実行

各 BC ごとに `go.mod` が分かれているため、テストは `cd <BC>` してから実行します。

```bash
# 全 BC のユニットテスト + 統合テスト
make test

# race detector 付き
make test-race

# 個別 BC を起動
make run-booking      # :8081
make run-catalog      # :8082
make run-identity     # :8083
```

## BC 間連携

- Booking BC は Catalog BC の `Lesson` を `shared/dto.LessonView` でのみ参照します
- Booking BC の `application.LessonQueryPort` interface を実装する側 (本サンプルでは stub) を `main.go` で組み立てます
- 実プロジェクトでは HTTP クライアントや message bus を実装します

## ライセンス

MIT License
