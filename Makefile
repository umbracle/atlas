SHELL := /bin/bash

ngrok:
	go build -o atlas cmd/main.go

protoc:
	protoc --go_out=. --go-grpc_out=. ./internal/proto/*.proto

