# 付録 C: Go 1.26 新機能と DDD への活用

『現場で使えるGo言語ドメイン駆動設計』付録 C のサンプルコード。

## 動作環境

- Go 1.26.0 以上 (本文サンプルとは別の go.mod)

このディレクトリは本文の `go.work` には含めません。Go 1.26 が必須のため、独立モジュールとして扱います。

## 実行

```bash
go test ./...
go test -race ./...
```

## 5 機能

- `errors_astype/`: `errors.AsType[T]` でドメインエラーを型安全に取り出す
- `self_ref_generics/`: ジェネリクス自己参照制約で通貨混在をコンパイル時に弾く
- `new_expr/`: `new(expr)` で Optional フィールドを 1 行で初期化
- `slog_multi/`: `slog.NewMultiHandler` で複数の出力先にログを配達
- `reflect_iter/`: `reflect.Type.Fields()` イテレータで集約の構造検査

## ライセンス

MIT License
