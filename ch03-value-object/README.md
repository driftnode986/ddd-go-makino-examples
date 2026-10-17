# 第3章 値オブジェクト

『現場で使える Go 言語ドメイン駆動設計』第3章のサンプルコード。

## 動作環境

- Go 1.25 または 1.26 の最新のパッチ版（`go.mod` の `go` の行は 1.25.0。標準ライブラリの脆弱性はパッチ版で直るため、`go version` で確かめて最新に上げておく）

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/email.go` — Email 値オブジェクト（コンストラクタで検証）
- `internal/booking/domain/money.go` — Money 値オブジェクト（int64 最小単位保持）
- `internal/booking/domain/hour.go` — Hour 値オブジェクト（営業時間の範囲制約）
- `internal/booking/domain/id.go` — 型安全な ID（ジェネリクス + phantom type）
- 各 `*_test.go` — 標準の `testing` パッケージだけで書いた Table-driven テスト

## トピック

- struct + 非公開フィールド + コンストラクタで不変条件を強制
- 値レシーバで「新しいインスタンスを返す」スタイル
- `==` で比較できない型を含む場合の `Equals` メソッド
- Java/Kotlin の `@Embeddable` / `data class` との対比

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
