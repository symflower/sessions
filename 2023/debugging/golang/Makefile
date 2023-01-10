reset:
	rm -f ./**/*_test.go go.mod go.sum
.PHONY: reset

init:
	go mod init debugging
	go get github.com/stretchr/testify/assert
.PHONY: init

example:
	go run main.go 8796
.PHONY: example

symflower:
	symflower --workspace hashing
.PHONY: symflower

test:
	go test ./hashing
.PHONY: test

all:
	make reset
	make init
	make example
	make symflower
	make test; true
	make reset
.PHONY: all
