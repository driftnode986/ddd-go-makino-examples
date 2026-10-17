# ddd-go-makino-examples

『[現場で使えるGo言語ドメイン駆動設計](https://www.amazon.co.jp/dp/B0H2R8N7WP)』のサンプルコード。

本に載せている Go のコードは、説明のための短い断片を除き、すべてこのリポジトリのファイルから取っています（本はこのリポジトリのコードの一部を抜粋して載せます）。各章のディレクトリで `go test ./...` が通ることを確かめながら読み進めてください。

## 動作環境

- Go 1.25 以降（1.25 / 1.26 系の最新のパッチ版を使ってください。付録 C だけは Go 1.26 が必須）
- Docker（PostgreSQL を使うテストを testcontainers で動かす場合のみ。第7章・第10章の Booking BC・第11章の `dbtest/`・付録 B）

## 構成

各章のディレクトリは独立した Go module です（リポジトリのルートには `go.mod` も `go.work` もありません）。章のディレクトリに移ってから実行します。

```
ch02-go-idioms/            第2章  Go の慣習と DDD の原則
ch03-value-object/         第3章  値オブジェクト
ch04-entity/               第4章  エンティティと不変条件の保護
ch05-domain-service/       第5章  ドメインサービス
ch06-aggregate-factory/    第6章  集約とファクトリ
ch07-repository/           第7章  リポジトリと永続化
ch08-application-service/  第8章  アプリケーションサービスとユースケース
ch09-layered-architecture/ 第9章  レイヤー化アーキテクチャ
ch10-bounded-context/      第10章 境界づけられたコンテキストとドメインイベント (go.work で 2 BC)
ch11-practical-project/    第11章 実プロジェクト構成例 (go.work で 3 BC + shared。dbtest は独立 module)
appendix-b-translation/    付録B  Java/Kotlin から Go への翻訳ガイドの Go のコード
appendix-c-go126/          付録C  Go 1.26 の新機能 (Go 1.26 が必須)
```

## 実行

```bash
# 単一 module の章 (第2〜9章・付録 B・付録 C)
(cd ch03-value-object && go test ./... && go test -race ./...)

# PostgreSQL を使うテスト (Docker が要る)
(cd ch07-repository && go test -tags=integration ./...)
(cd appendix-b-translation && go test ./...)

# go.work で複数 module をまとめた章 (第10章・第11章)
(cd ch10-bounded-context && go test ./booking/... ./notification/...)   # booking/db は Docker が要る

(cd ch11-practical-project && make test)      # 各 BC のディレクトリで go test ./...
(cd ch11-practical-project && make test-db)   # dbtest (GOWORK=off、Docker が要る)
```

`go.work` のあるディレクトリには `go.mod` が無いので、そこで `go test ./...` を実行すると失敗します。第10章のように module ごとのパターンを並べるか、BC のディレクトリに移ってから実行してください。

章ごとの詳しい構成・実行方法は各ディレクトリの `README.md` を参照してください。

## ライセンス

MIT License（[LICENSE](./LICENSE) 参照）
