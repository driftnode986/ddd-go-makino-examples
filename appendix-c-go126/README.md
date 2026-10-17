# 付録C Go 1.26 新機能と DDD への活用

『現場で使えるGo言語ドメイン駆動設計』付録 C のサンプルコード。

## 動作環境

- **Go 1.26 が必須**（本文サンプルとは別の `go.mod`）。1.26 系の最新のパッチ版を使ってください（古いパッチ版の標準ライブラリには、修正済みの脆弱性が残っています）

このディレクトリは、ほかの章と同じく独立した module です（どの `go.work` にも入れていません）。本文の章ディレクトリは Go 1.25 のままビルドできます。

## 実行

```bash
go test ./...
go test -race ./...
```

## 構成（Go 1.26 新機能 5 件）

- `errors_astype/` — `errors.AsType[T]` でドメインエラーを型安全に取り出す
- `self_ref_generics/` — ジェネリクス自己参照制約で通貨混在をコンパイル時に弾く（`Money[JPY]` phantom type）
- `new_expr/` — `new(expr)` で Optional フィールドを 1 行で初期化
- `slog_multi/` — `slog.NewMultiHandler` で複数の出力先にログを配達
- `reflect_iter/` — `reflect.Type.Fields()` イテレータで集約の構造検査

## ライセンス

MIT License（ルート [LICENSE](../LICENSE) 参照）
