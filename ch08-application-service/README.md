# 第8章 アプリケーションサービスとユースケース

『現場で使えるGo言語ドメイン駆動設計』第8章のサンプルコード。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください）

## 実行

```bash
go test ./...                       # ユニットテスト
go test -race ./...                 # race detector 付き
go build ./...                      # ビルドのみ確認
go run ./cmd/api                    # サーバ起動 (:8080。ハンドラの登録は第9章のため、この章ではどのパスも 404)
```

## 構成

- `internal/order/domain/` — ドメイン層（Order 集約 + 値オブジェクト + Repository interface）
- `internal/order/application/` — アプリケーション層（PlaceOrderUseCase + ListOrdersByCustomerHandler + DTO + ApplicationError）
- `internal/order/infrastructure/` — インフラ層（in-memory リポジトリ + Stub Pricing）
- `cmd/api/main.go` — main 関数で依存を組み立てる例（DI コンテナ不使用）

## トピック

- 1 ユースケース = 1 メソッドの構造（PlaceOrderUseCase は NewOrder → AddItem → Confirm → Add の順に呼ぶ）
- ドメイン層を呼び出すオーケストレータとしてのアプリケーションサービス（業務ルールは書かない）
- Command（書き込み）と Query（読み取り）の分離（CQRS の入り口）
- DTO とドメインモデルの詰め替え
- 3 段エラーマッピング（domain error → ApplicationError → HTTP status）
- `context.Context` での request_id 伝播
- `log/slog` の `InfoContext` 活用
- main 関数での手書きワイヤリング（Java/Kotlin の `@Service` + DI コンテナとの対比）

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
