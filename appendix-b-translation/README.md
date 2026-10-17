# 付録B Java/Kotlin DDD から Go DDD への翻訳ガイド

付録 B の Go のコードを、動く形でまとめた独立モジュールです（本文の章とは別の `go.mod`）。

- `internal/identity/` — 会員（`User`・`Email`）、`UserRepository` とその PostgreSQL 実装（sqlc）、`RegisterUserUseCase`（`TxManager` で 1 トランザクションにまとめる）
- `internal/booking/` — 注文（`Order`）、`PricingService`、`PlaceOrderUseCase`（DI の例）、`application/publishing/`（イベントを発行する例）、エラーから HTTP ステータスへの対応付け
- `internal/platform/` — 本番用の時計（`system`）、テスト用の固定の時計（`test`）、pgx のトランザクションを context で渡す `pgxtx`
- `db/` — テーブル定義と sqlc のクエリ。`sqlc.yaml` から `sqlc generate` で `sqlcgen/` を作る（生成済みのコードもコミットしている）

## 実行

```bash
go test ./...      # PostgreSQL を使うテストは testcontainers で動かすので Docker が要る
```

`cmd/api` は `DATABASE_URL` の PostgreSQL（`db/booking_schema.sql` を流したもの）に接続して `:8080` で起動します。

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
