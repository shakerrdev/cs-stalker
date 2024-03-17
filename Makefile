.DEFAULT_GOAL := run

run: cmd/cs-stalker/main.go
	go run cmd/cs-stalker/main.go
