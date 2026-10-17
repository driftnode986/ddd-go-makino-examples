// Package db はテーブル定義を埋め込む。PostgreSQL を使うテストが流す。
package db

import _ "embed"

// IdentitySchema は identity_schema.sql（users テーブル）の内容。
//
//go:embed identity_schema.sql
var IdentitySchema string

// BookingSchema は booking_schema.sql（orders テーブル）の内容。
//
//go:embed booking_schema.sql
var BookingSchema string
