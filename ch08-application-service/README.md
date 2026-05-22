# 第8章 アプリケーションサービスとユースケース

『現場で使える Go 言語ドメイン駆動設計』第8章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 構成

- `internal/order/domain/` — ドメイン層 (Order 集約 + 値オブジェクト + リポジトリ interface)
- `internal/order/application/` — アプリケーション層 (PlaceOrderUseCase + ListOrdersByCustomerHandler + DTO + ApplicationError)
- `internal/order/infrastructure/` — インフラ層 (in-memory リポジトリ + Stub Pricing)
- `cmd/api/main.go` — main 関数で依存を組み立てる例

## 実行

```bash
go test ./...                       # ユニットテスト
go test -race ./...                 # race detector 付き
go build ./...                      # ビルドのみ確認
```

## ライセンス

MIT License
