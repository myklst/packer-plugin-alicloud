NAME=alicloud
BINARY=packer-plugin-${NAME}

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module \s*//')

.PHONY: dev

build:
	@go build -o ${BINARY}

dev:
	go build -ldflags="-X '${PLUGIN_FQN}/version.VersionPrerelease=dev'" -o ${BINARY}
	packer plugins install --path ${BINARY} "$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

test:
	@go test -race -count $(COUNT) $(TEST) -timeout=3m

install-packer-sdc: ## Install packer sofware development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BINARY}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

generate: install-packer-sdc
	@go generate ./...
	@rm -rf .docs
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst .docs/
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "hashicorp"
	@rm -r ".docs"

.PHONY: packer-plugin-install
packer-plugin-install:
	go build -ldflags="-X github.com/myklst/packer-plugin-alicloud/version.VersionPrerelease=dev" -o packer-plugin-alicloud
	packer plugins install --path packer-plugin-alicloud github.com/myklst/alicloud

.PHONY: go-test
go-test:
	export PACKER_ACC=1
	go test ./...

.PHONY: go-generate
go-generate:
	go generate ./...
