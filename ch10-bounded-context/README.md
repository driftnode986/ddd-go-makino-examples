# 第10章: 境界づけられたコンテキストとドメインイベント

『現場で使えるGo言語ドメイン駆動設計』第10章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 構成

```text
ch10-bounded-context/
  go.work                                 # 2 BC をワークスペースで束ねる
  booking/
    go.mod                                # Booking BC の独立 module
    internal/booking/
      domain/                             # ドメイン層 (Order 集約 + OrderPlaced イベント)
      application/                        # ユースケース + OutboxPort
      infrastructure/                     # in-memory リポジトリ + in-memory outbox
  notification/
    go.mod                                # Notification BC の独立 module
    internal/notification/
      events/                             # 公表された言語 (OrderPlaced DTO)
      application/                        # OrderPlacedHandler + IdempotencyStore
      infrastructure/                     # in-memory mailer
```

## 実行

```bash
go test ./booking/... ./notification/...                # ユニットテスト
go test -race ./booking/... ./notification/...          # race detector
```

## ライセンス

MIT License
