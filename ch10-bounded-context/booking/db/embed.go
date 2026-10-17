// Package db は outbox_events テーブルの定義 (outbox_events.sql) を埋め込む。
// 本章の repo は in-memory の outbox で動かすため、PostgreSQL 版の outbox は
// テーブル定義だけを置き、統合テストで実際の PostgreSQL に流して確かめる。
package db

import _ "embed"

// OutboxSchema は outbox_events.sql の内容。
//
//go:embed outbox_events.sql
var OutboxSchema string
