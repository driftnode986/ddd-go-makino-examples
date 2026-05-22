# 第3章 値オブジェクト

『現場で使えるGo言語ドメイン駆動設計』第3章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/email.go` — Email 値オブジェクト
- `internal/booking/domain/money.go` — Money 値オブジェクト
- `internal/booking/domain/hour.go` — Hour 値オブジェクト
- `internal/booking/domain/id.go` — 型安全な ID（ジェネリクス）
- 各 `*_test.go` — Table-driven テスト + go-cmp

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
