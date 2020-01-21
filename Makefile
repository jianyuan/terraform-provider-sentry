export CGO_ENABLED = 0
VERSION = $(shell git describe --tags --match='v*')

CROSSBUILD_OS   = linux darwin
CROSSBUILD_ARCH = 386 amd64
OSARCH_COMBOS   = $(foreach os,$(CROSSBUILD_OS),$(addprefix $(os)_,$(CROSSBUILD_ARCH)))

default: build

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

vet:
	@echo ">> vetting code"
	@go vet ./...

test:
	@echo ">> testing code"
	@TF_ACC=1 go test -race -v ./...

deps:
	@echo ">> downloading deps"
	@go mod download

build:
	@echo ">> building binaries"
	@go build -o terraform-provider-sentry

crossbuild: $(GOPATH)/bin/gox
	@echo ">> cross-building"
	@gox -arch="$(CROSSBUILD_ARCH)" -os="$(CROSSBUILD_OS)" -output="binaries/{{.OS}}_{{.Arch}}/terraform-provider-sentry_$(VERSION)"

$(GOPATH)/bin/gox:
	# Need to disable modules for this to not pollute go.mod
	@GO111MODULE=off go get -u github.com/mitchellh/gox

release: crossbuild bin/github-release
	@echo ">> uploading release ${VERSION}"
	@mkdir -p releases
	@set -e; for OSARCH in $(OSARCH_COMBOS); do \
		zip -j releases/terraform-provider-sentry_$(VERSION)_$$OSARCH.zip binaries/$$OSARCH/terraform-provider-sentry_* > /dev/null; \
		./bin/github-release upload -t ${VERSION} -n terraform-provider-sentry_$(VERSION)_$$OSARCH.zip -f releases/terraform-provider-sentry_$(VERSION)_$$OSARCH.zip; \
	done

bin/github-release:
	@mkdir -p bin
	@curl -sL 'https://github.com/aktau/github-release/releases/download/v0.7.2/linux-amd64-github-release.tar.bz2' | tar xjf - --strip-components 3 -C bin

.PHONY: all style vet test deps build crossbuild release
