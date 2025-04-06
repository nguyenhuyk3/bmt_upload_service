SERVICE_NAME=bmt_product

DB_URL=postgres://postgres:anhiuemlove33@127.0.0.1:5432/bmt_product?sslmode=disable

run:
	go run .\cmd\server\main.go

wire:
	wire ./internal/injectors/

.PHONY: run slqc wire
