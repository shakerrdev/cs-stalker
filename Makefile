.DEFAULT_GOAL := all

.PHONY: all
all: test run

.PHONY: test
# grep to exclude application logs
test:
	go test -v ./... | grep -v "^[0-9]\{4\}/[0-9]\{2\}/[0-9]\{2\} [0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}"

.PHONY: run
run:
	go run ./cmd/cs-stalker
