.PHONE: all
all: test

.PHONY: deps
deps:
	@go mod download

.PHONY: test
test:
	@TF_ACC=1 go test -race -v ./...
