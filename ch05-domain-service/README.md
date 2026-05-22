# 第5章 ドメインサービス

『現場で使えるGo言語ドメイン駆動設計』第5章のサンプルコード。

## 動作環境

- Go 1.25.0+

## 実行

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 構成

- `internal/booking/domain/` — Clock, Booking, ConflictDetector, ExpirationChecker
- `internal/pricing/domain/` — Money, DiscountRate, CalculateDiscountedPrice (関数版ドメインサービス)
- `internal/coupon/domain/` — CouponCode, RandomSource, CodeGenerator

## トピック

- ドメインサービスの配置判断 (6 軸)
- 関数 vs 構造体の使い分け
- Clock interface で時刻依存を分離
- RandomSource interface で乱数依存を分離
- 戻り値は値オブジェクトで返す
- 小さい interface (`BookingFinder`) でテストダブルを簡素化

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
