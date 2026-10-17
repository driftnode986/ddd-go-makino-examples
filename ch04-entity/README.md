# 第4章 エンティティと不変条件保護

『現場で使えるGo言語ドメイン駆動設計』第4章のサンプルコード。

## 動作環境

- Go 1.25 または 1.26 の最新のパッチ版（`go.mod` の `go` の行は 1.25.0。標準ライブラリの脆弱性はパッチ版で直るため、`go version` で確かめて最新に上げておく）

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/id.go` — phantom type で型安全な ID（`BookingTag` / `MemberTag` / `LessonTag` のマーカー型サフィックス）
- `internal/booking/domain/booking.go` — Booking エンティティ（Pending → Confirmed → Completed / Cancelled の状態遷移）
- `internal/booking/domain/booking_test.go` — 状態遷移と不変条件のテスト
- `internal/identity/domain/user.go` — User エンティティ（登録 / 退会）
- `internal/identity/domain/user_test.go` — User のテスト

## トピック

- ID 同一性とライフサイクルを持つドメインオブジェクト
- ポインタレシーバでの状態変更メソッド
- 不変条件をエンティティ本体で守る（公開フィールドで状態変更を許さない）
- 楽観的ロックの `version` フィールド（実装の深掘りは Ch7）
- マーカー型名衝突の回避（4.5 節）

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
