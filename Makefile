GO_VER        ?= go
SWEEP         ?= cloud
SWEEP_TIMEOUT ?= 360m

default: testacc

.PHONY: deps
deps:
	$(GO_VER) mod download

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 $(GO_VER) test ./... -v $(TESTARGS) -timeout 120m

sweep: ## Run sweepers
	# make sweep SWEEPARGS=-sweep-run=sentry_team
	# set SWEEPARGS=-sweep-allow-failures to continue after first failure
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	$(GO_VER) test ./... -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout $(SWEEP_TIMEOUT)

sweeper: ## Run sweepers with failures allowed
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	$(GO_VER) test ./... -v -tags=sweep -sweep=$(SWEEP) -sweep-allow-failures -timeout $(SWEEP_TIMEOUT)