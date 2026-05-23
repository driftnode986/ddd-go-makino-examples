# ddd-go-makino-examples

『現場で使える Go 言語ドメイン駆動設計』（牧野 誠）のサンプルコード。

書籍内のすべての Go コード片はこのリポジトリで動作確認できます。`go test ./...` で各章のテストが PASS することを確認しながら読み進めてください。

## 動作環境

- Go 1.25.0 以上（付録 C のみ Go 1.26.0 以上）
- PostgreSQL 16 以上（Ch07 / Ch11 の統合テストを動かす場合のみ）
- Docker（Testcontainers / docker compose を使う場合のみ）

## 構成

各章は独立した Go モジュールです。章ディレクトリ単位で `go test ./...` が走ります。

```
ch03-value-object/         第3章  値オブジェクト
ch04-entity/               第4章  エンティティと不変条件保護
ch05-domain-service/       第5章  ドメインサービス
ch06-aggregate-factory/    第6章  集約とファクトリ
ch07-repository/           第7章  リポジトリと永続化
ch08-application-service/  第8章  アプリケーションサービスとユースケース
ch09-layered-architecture/ 第9章  レイヤー化アーキテクチャ
ch10-bounded-context/      第10章 境界づけられたコンテキストとドメインイベント (go.work で 2 BC)
ch11-practical-project/    第11章 実プロジェクト構成例 (go.work で 3 BC + shared)
appendix-c-go126/          付録C  Go 1.26 新機能 (独立 go.mod、Go 1.26 必須)
```

付録 B（Java/Kotlin → Go 翻訳ガイド）は本文 Ch3-Ch11 のコードを再掲する構成のため、専用ディレクトリは用意していません。各章のコードがそのまま参考になります。

## 実行

```bash
# 単一章 (Ch3-Ch9, 付録 C)
cd ch03-value-object
go test ./...
go test -race ./...
go test -cover ./...

# 複数 module 構成 (Ch10, Ch11)
cd ch10-bounded-context
go test ./booking/... ./notification/...

cd ch11-practical-project
make test
```

章ごとの詳しい構成・実行方法は各ディレクトリの `README.md` を参照してください。

## ライセンス

MIT License（[LICENSE](./LICENSE) 参照）
