# ch11-practical-project Makefile
# トレーニング予約 SaaS の 3 BC モノレポを操作する最小コマンド集です。

.PHONY: up down test test-race lint build run-booking run-catalog run-identity

up:
	docker compose up -d

down:
	docker compose down

test:
	cd booking && go test ./...
	cd catalog && go test ./...
	cd identity && go test ./...

test-race:
	cd booking && go test -race ./...
	cd catalog && go test -race ./...
	cd identity && go test -race ./...

build:
	cd booking && go build ./...
	cd catalog && go build ./...
	cd identity && go build ./...

run-booking:
	cd booking && go run ./cmd/api

run-catalog:
	cd catalog && go run ./cmd/api

run-identity:
	cd identity && go run ./cmd/api
