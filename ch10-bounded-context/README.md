# 第10章 境界づけられたコンテキストとドメインイベント

『現場で使える Go 言語ドメイン駆動設計』第10章のサンプルコード。

## 動作環境

- Go 1.25.0 以上

## 実行

```bash
# 両 BC まとめてテスト (go.work モード)
go test ./booking/... ./notification/...
go test -race ./booking/... ./notification/...

# 個別の BC を単独でテスト
cd booking
go test ./...

cd ../notification
go test ./...
```

## 構成

```text
ch10-bounded-context/
  go.work                                 # 2 BC をワークスペースで束ねる
  booking/
    go.mod                                # Booking BC の独立 module
    internal/booking/
      domain/                             # ドメイン層 (Order 集約 + OrderPlaced ドメインイベント)
      application/                        # ユースケース + OutboxPort
      infrastructure/                     # in-memory リポジトリ + in-memory outbox
  notification/
    go.mod                                # Notification BC の独立 module
    internal/notification/
      events/                             # 公表された言語 (OrderPlaced DTO)
      application/                        # OrderPlacedHandler + IdempotencyStore
      infrastructure/                     # in-memory mailer
```

Booking BC は Notification BC のコードを import しません。Notification BC は `events/OrderPlaced` DTO のみを公表された言語として持ち、Booking 側のドメインモデルに依存しません。

## トピック

- 境界づけられたコンテキストを `go.mod` 単位で物理分離
- `go.work` で複数 BC を 1 リポジトリにまとめる
- ドメインイベント（`OrderPlaced`）の発行
- アウトボックスパターン（書き込みトランザクションと整合する形でイベントを永続化）
- 冪等性ストアでの重複処理ガード
- 公表された言語（Published Language）としての DTO 設計

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
