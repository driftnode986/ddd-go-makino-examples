# 第4章 エンティティと不変条件保護

『現場で使えるGo言語ドメイン駆動設計』第4章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/id.go` — phantom type 型安全 ID
- `internal/booking/domain/booking.go` — Booking エンティティ (予約/キャンセル/完了)
- `internal/booking/domain/booking_test.go` — Booking テスト
- `internal/identity/domain/user.go` — User エンティティ (登録/退会)
- `internal/identity/domain/user_test.go` — User テスト

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
