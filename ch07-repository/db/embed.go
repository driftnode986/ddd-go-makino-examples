// Package db は bookings テーブルの定義 (schema.sql) をバイナリに埋め込む。
// 統合テストの MigrateSchema が読み込む。
package db

import _ "embed"

// Schema は schema.sql の内容。
//
//go:embed schema.sql
var Schema string
