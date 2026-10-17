# ch11-practical-project Makefile
# トレーニング予約 SaaS の 3 BC モノレポを操作する最小コマンド集です。

# #region minimal
.PHONY: up down test build run-booking run-catalog run-identity

up:
	docker compose up -d

down:
	docker compose down

test:
	cd booking && go test ./...
	cd catalog && go test ./...
	cd identity && go test ./...

build:
	cd booking && go build ./...
	cd catalog && go build ./...
	cd identity && go build ./...

run-booking:
	cd booking && go run ./cmd/api
# #endregion

run-catalog:
	cd catalog && go run ./cmd/api

run-identity:
	cd identity && go run ./cmd/api

.PHONY: test-race test-db

test-race:
	cd booking && go test -race ./...
	cd catalog && go test -race ./...
	cd identity && go test -race ./...

# SQL（各 BC のマイグレーションと booking の queries.sql）を PostgreSQL に流して確かめる。
# Docker が要る。dbtest は go.work に入れない独立 module なので GOWORK=off で動かす。
test-db:
	cd dbtest && GOWORK=off go test ./...
