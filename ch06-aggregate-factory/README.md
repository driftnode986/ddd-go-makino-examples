# 第6章 集約とファクトリ

『現場で使えるGo言語ドメイン駆動設計』第6章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/order/domain/id.go` — OrderID / OrderItemID / CustomerID / ProductID (phantom type)
- `internal/order/domain/clock.go` — Clock interface (SystemClock + FixedClock)
- `internal/order/domain/money.go` — Money / DiscountRate 値オブジェクト
- `internal/order/domain/order_item.go` — OrderItem 集約内部の値オブジェクト
- `internal/order/domain/order.go` — Order 集約ルート
- `internal/order/domain/factory.go` — NewOrder / OrderFactory / Option 関数 / ReconstructOrder
- `internal/order/domain/order_test.go` — 集約の不変条件・defensive copy テスト
- `internal/order/domain/factory_test.go` — ファクトリのテスト

## トピック

- Vernon の Effective Aggregate Design 4 原則 + 本書独自の 1 原則
- 集約ルートのメソッドで内部を変更
- ID 参照 (他集約をオブジェクト参照しない)
- slices.Clone で defensive copy
- ファクトリ関数 vs ファクトリ構造体
- Option 関数パターン
- 再構築用ファクトリ (ReconstructOrder)

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
