SWEEP ?= cloud
SWEEP_TIMEOUT ?= 360m

default: fmt lint install generate

.PHONY: build
build:
	go build -v ./...

.PHONY: install
install: build
	go install -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./internal/apiclient
	go generate ./internal/sentrydata
	go generate ./internal/providergen
	go generate ./

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: test
test:
	go test ./... -v -cover -timeout=120s -parallel=10 $(TESTARGS)

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v -cover -timeout 120m $(TESTARGS)

.PHONY: sweep
sweep:
	# make sweep SWEEPARGS=-sweep-run=sentry_team
	# set SWEEPARGS=-sweep-allow-failures to continue after first failure
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./... -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout $(SWEEP_TIMEOUT)

.PHONY: sweeper
sweeper:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./... -v -tags=sweep -sweep=$(SWEEP) -sweep-allow-failures -timeout $(SWEEP_TIMEOUT)
