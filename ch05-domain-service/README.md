# 第5章 ドメインサービス

『現場で使えるGo言語ドメイン駆動設計』第5章のサンプルコード。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください）

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/` — Clock interface、Booking、Hour・HoursOverlap、ConflictDetector、ExpirationChecker
- `internal/pricing/domain/` — Money、DiscountRate、CalculateDiscountedPrice（関数版ドメインサービス）
- `internal/coupon/domain/` — CouponCode、RandomSource interface、CodeGenerator

## トピック

- ドメインサービスの配置判断（6 軸）
- 関数版ドメインサービス vs 構造体版ドメインサービスの使い分け
- Clock interface で時刻依存を分離（SystemClock + FixedClock）
- RandomSource interface で乱数依存を分離（テスト時は決定的なソース）
- 戻り値は値オブジェクトで返す（プリミティブを返さない）
- 小さい interface（`BookingFinder`）でテストダブルを簡素化

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
