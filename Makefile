IMAGE_COPY_REGIONS ?= []

.PHONY: packer-plugin-install
packer-plugin-install:
	go build -ldflags="-X github.com/myklst/packer-plugin-alicloud/version.VersionPrerelease=dev" -o packer-plugin-alicloud
	packer plugins install --path packer-plugin-alicloud github.com/myklst/alicloud

.PHONY: go-test
go-test:
	go test ./...
