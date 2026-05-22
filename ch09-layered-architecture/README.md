# 第9章: レイヤー化アーキテクチャとパッケージ分割

『現場で使えるGo言語ドメイン駆動設計』第9章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 構成

```text
ch09-layered-architecture/
  cmd/api/main.go              # 全レイヤーのワイヤリング + graceful shutdown
  internal/
    order/
      domain/                  # ドメイン層 (エンティティ・値オブジェクト・リポジトリ interface)
      application/             # アプリケーション層 (ユースケース・クエリ・ApplicationError)
      infrastructure/          # インフラ層 (in-memory リポジトリ)
      presentation/            # プレゼンテーション層 (HTTP ハンドラ・ミドルウェア・DTO)
```

依存方向は内向き一方向: Presentation → Application → Domain ← Infrastructure。

## 実行

```bash
go test ./...                       # ユニットテスト
go test -race ./...                 # race detector
go test -cover ./...                # カバレッジ
go build ./cmd/api                  # API サーバビルド
go run ./cmd/api                    # API サーバ起動 (:8080)
```

## API エンドポイント

- `POST /orders` — 注文作成
- `GET /orders/{id}` — 注文詳細取得
- `GET /orders?customer_id=cus_001` — 顧客別注文一覧

## ライセンス

MIT License
