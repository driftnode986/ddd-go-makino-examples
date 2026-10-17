# 第2章 Go の慣習と DDD の原則を擦り合わせる

『現場で使える Go 言語ドメイン駆動設計』第2章のサンプルコード。

## 動作環境

- Go 1.25.0 以上

## 実行

```bash
cd ch02-go-idioms
go test ./...
```

## 構成

- `internal/` — 章の完成形。本文の各節のコードが 1 つのモジュールで動く
  - `booking/domain/` — Booking 集約・リポジトリの interface・ドメインエラー・時間帯の衝突判定
  - `booking/application/` — レッスン予約のユースケース（log/slog で記録する版）と、iter で予約を 1 件ずつ処理する例
  - `booking/adapter/http/` — DTO と、ドメインエラーを HTTP ステータスに対応させるハンドラ
  - `booking/adapter/postgres/` — Persistence Model（`db:` タグ付き）とマッピング関数
  - `identity/domain/` — Email・Money・User
  - `pricing/domain/` — 時計を注入する PricingService
  - `shared/domain/` — ジェネリクスの型安全な ID
  - `before/` `after/` — 2.15 節の演習（リファクタ前とリファクタ後）
- `steps/` — 本文が同じファイルを別の形で示している節の版。完成形と混ざらないよう、節ごとに別のパッケージにしてある
  - `s21-small-interface/` — 2.1 節（3 メソッドのリポジトリと、中身を省略したメソッドの多い Booking）
  - `s22-postgres-model/` — 2.2 節（タグなしの Persistence Model）
  - `s23-errors/` — 2.3 節（errors.As で衝突を判別するユースケース）
  - `s25-receivers/` — 2.5 節（ID と状態だけに絞った Booking でポインタレシーバを揃える）
  - `s25-mixed-receivers/` — 2.5 節のアンチパターン（レシーバの混在）
  - `s27-behavior-naming/` — 2.7 節（Setter とビヘイビア指向のメソッド）
