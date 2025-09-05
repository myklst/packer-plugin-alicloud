BINARY=packer-plugin-alicloud
PLUGIN_FQN="$(shell grep -E '^module' <go.mod | sed -E 's/module *//')"

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

.PHONY: dev

build:
	@go build -o ${BINARY}

dev: clean
	@go build -ldflags="-X '${PLUGIN_FQN}/version.VersionPrerelease=dev'" -o '${BINARY}'
	@packer plugins install --path ${BINARY} "$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

docs: install-packer-sdc
	@rm -rf .docs docs-partials .web-docs/components
	@go generate ./...
	@$(shell go env GOPATH)/bin/packer-sdc renderdocs -src ./docs -partials ./docs-partials -dst ./.docs
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "hashicorp"
	@rm -r ".docs"

test: dev
	@go test -race -count $(COUNT) $(TEST) -timeout=5m

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

testacc-builder-eds: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v ./builder/eds/builder_acc_test.go -timeout=120m

testacc-datasource-images: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v ./datasource/ecsimage/data_acc_test.go -timeout=120m

# Install packer sofware development command
install-packer-sdc:
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

plugin-check: install-packer-sdc build
	@$(shell go env GOPATH)/bin/packer-sdc plugin-check ${BINARY}

.IGNORE:
clean:
	@packer plugins remove "github.com/myklst/alicloud"
	@rm -rf packer-plugin-alicloud
