// Package dbtest は、本章の SQL（各 BC のマイグレーションと booking の queries.sql）を
// 実際の PostgreSQL に流して確かめるための module。
//
// 本章のサンプルは in-memory リポジトリで動かすため、BC の module には PostgreSQL の
// ドライバを入れていない。SQL の確認だけをこの独立 module で行う（go.work には入れない）。
//
//	cd dbtest && GOWORK=off go test ./...
package dbtest
