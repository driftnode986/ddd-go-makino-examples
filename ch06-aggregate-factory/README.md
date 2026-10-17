# 第6章 集約とファクトリ

『現場で使える Go 言語ドメイン駆動設計』第6章のサンプルコード。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください）

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/order/domain/id.go` — OrderID / CustomerID / ProductID（phantom type）
- `internal/order/domain/clock.go` — Clock interface（SystemClock + FixedClock）
- `internal/order/domain/money.go` — Money / DiscountRate 値オブジェクト
- `internal/order/domain/order_item.go` — OrderItem（集約内部の値オブジェクト）
- `internal/order/domain/order.go` — Order 集約ルート
- `internal/order/domain/factory.go` — NewOrder / OrderFactory / Option 関数 / ReconstructOrder
- `internal/order/domain/order_test.go` — 集約の不変条件と defensive copy のテスト
- `internal/order/domain/factory_test.go` — ファクトリのテスト
- `internal/order/domain/factory_example_test.go` — 6.9 の Option 関数の呼び出し例（Example 関数）

## トピック

- Vernon の Effective Aggregate Design の 4 原則と、本書が加える 2 つの観点（振る舞いを集約ルートのメソッドに置く・`slices.Clone` で内部スライスを返す）
- 集約ルートのメソッドで内部を変更する（公開フィールドを置かない）
- ID 参照（他集約をオブジェクト参照しない）
- `slices.Clone` で defensive copy
- ファクトリ関数 vs ファクトリ構造体
- Option 関数パターン（Java/Kotlin の Builder 相当）
- 再構築用ファクトリ（`ReconstructOrder`、永続化からの復元）

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
