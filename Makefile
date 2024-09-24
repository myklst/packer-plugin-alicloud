.PHONY: packer-plugin-install
packer-plugin-install:
	go generate ./...
	rm -rf /home/user/.config/packer/plugins/github.com/myklst/alicloud/* && \
	go build -ldflags="-X github.com/myklst/packer-plugin-alicloud/version.VersionPrerelease=dev" -o packer-plugin-alicloud
	packer plugins install --path packer-plugin-alicloud github.com/myklst/alicloud
	go test ./...
