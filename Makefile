BUILDDIR=./build
GOTIFY_VERSION=master
PLUGIN_NAME=broadcasts
LICENSE_DIR=./OPENSOURCELICENSE

download-tools:
	GO111MODULE=off go get -u github.com/gotify/plugin-api/cmd/gomod-cap

create-build-dir:
	mkdir -p ${BUILDDIR} || true

check-go:
	go test ./...

update-go-mod: create-build-dir
	wget -LO ${BUILDDIR}/gotify-server.mod https://raw.githubusercontent.com/gotify/server/${GOTIFY_VERSION}/go.mod
	gomod-cap -from ${BUILDDIR}/gotify-server.mod -to go.mod
	rm ${BUILDDIR}/gotify-server.mod || true
	go mod tidy

check-go-mod: create-build-dir
	wget -LO ${BUILDDIR}/gotify-server.mod https://raw.githubusercontent.com/gotify/server/${GOTIFY_VERSION}/go.mod
	gomod-cap -from ${BUILDDIR}/gotify-server.mod -to go.mod -check=true
	rm ${BUILDDIR}/gotify-server.mod || true

build: create-build-dir update-go-mod
	CGO_ENABLED=1 go build -o build/${PLUGIN_NAME}-$$(go env -json | jq -r ".GOOS")-$$(go env -json | jq -r ".GOARCH")${GOARM}-for-gotify-${GOTIFY_VERSION}.so -buildmode=plugin

extract-licenses:
	go mod vendor
	mkdir ${LICENSE_DIR} || true
	for LICENSE in $(shell find vendor/* -name LICENSE); do \
		DIR=`echo $$LICENSE | tr "/" _ | sed -e 's/vendor_//; s/_LICENSE//'` ; \
        cp $$LICENSE ${LICENSE_DIR}/$$DIR ; \
	done


check: check-go check-go-mod

.PHONY: build download-tools check-go check-go-mod update-go-mod
