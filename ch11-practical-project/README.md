# 第11章 実プロジェクト構成例とテスト戦略

『現場で使える Go 言語ドメイン駆動設計』第11章のサンプルコード。

トレーニング予約 SaaS を 3 つの境界づけられたコンテキスト（Booking / Catalog / Identity）で構成し、`go.work` で 1 リポジトリにまとめる capstone 章のサンプルです。

## 動作環境

- Go 1.25.0 以上
- Docker（オプション、PostgreSQL を起動する場合のみ）

## 実行

各 BC ごとに `go.mod` が分かれているため、Makefile 経由で実行することを推奨します。

```bash
make test          # 全 BC のユニットテスト
make test-race     # race detector 付き
make build         # 全 BC をビルド
make up            # PostgreSQL を docker compose で起動
make down          # PostgreSQL を停止
make run-booking   # Booking BC を起動 (:8081)
make run-catalog   # Catalog BC を起動 (:8082)
make run-identity  # Identity BC を起動 (:8083)
```

個別 BC を直接実行する場合は `cd booking && go test ./...` のように章ディレクトリに移動してから実行します。

## 構成

```text
ch11-practical-project/
  go.work                       # 4 module をまとめる
  docker-compose.yaml           # PostgreSQL 17
  Makefile                      # 実行コマンド一覧
  shared/                       # 公表された言語 DTO
    dto/lesson_view.go
  booking/                      # Booking BC (フル実装サンプル)
    cmd/api/main.go
    internal/booking/
      domain/                   # Booking 集約 + Repository interface
      application/              # ユースケース + クエリ
      infrastructure/memory/    # in-memory リポジトリ
      presentation/http/        # chi ハンドラ + httptest
      testdata/                 # Test Data Builder
  catalog/                      # Catalog BC (Lesson 集約のみ、最小構成)
  identity/                     # Identity BC (Member 集約のみ、最小構成)
```

## BC 間連携

- Booking BC は Catalog BC の `Lesson` を `shared/dto.LessonView` でのみ参照します
- Booking BC の `application.LessonQueryPort` interface を実装する側（本サンプルでは stub）を `main.go` で組み立てます
- 実プロジェクトでは HTTP クライアントや message bus で実装します

## トピック

- 3 BC × Domain / Application / Infrastructure / Presentation の 4 層構成
- `go.work` + `replace` ディレクティブで shared パッケージを共有
- テストピラミッド（ユニット中心、統合テストは Testcontainers）
- mock を使わないテスト戦略（Test Data Builder + Stub）
- graceful shutdown（`context.WithCancel` + `http.Server.Shutdown`）
- CI 構成例（`working-directory: <BC>` で BC ごとにジョブを分ける）

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
