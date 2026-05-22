# ddd-go-makino-examples

『現場で使えるGo言語ドメイン駆動設計』（牧野 誠）のサンプルコード。

## 動作環境

- Go 1.25.0+（付録 C のみ 1.26.0+）
- PostgreSQL 16+（Ch07 以降の Testcontainers 統合テスト用）
- Docker（Testcontainers 用、ローカル統合テストを走らせる場合）

## 構成

各章はモジュール独立。章ディレクトリ単位で `go test ./...` 可能。

```
ch03-value-object/         第3章 値オブジェクト
ch04-entity/               第4章 エンティティ
ch05-domain-service/       第5章 ドメインサービス
ch06-aggregate-factory/    第6章 集約とファクトリ
ch07-repository/           第7章 リポジトリと永続化
ch08-application-service/  第8章 アプリケーションサービス
ch09-layered-architecture/ 第9章 レイヤー化アーキテクチャ
ch10-bounded-context/      第10章 境界づけられたコンテキスト
ch11-practical-project/    第11章 実プロジェクト構成例
appendix-b-translation/    付録B Java/Kotlin → Go 翻訳
appendix-c-go126/          付録C Go 1.26 新機能
```

## 実行

```bash
cd ch03-value-object
go test ./...
go test -race ./...
go test -cover ./...
```

## ライセンス

MIT License（[LICENSE](./LICENSE) 参照）
