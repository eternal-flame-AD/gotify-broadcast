BUILDDIR=./build
GOTIFY_VERSION=master
PLUGIN_NAME=broadcasts
PLUGIN_ENTRY=plugin.go
GO_VERSION=`cat $(BUILDDIR)/gotify-server-go-version`

download-tools:
	GO111MODULE=off go get -u github.com/gotify/plugin-api/cmd/gomod-cap

create-build-dir:
	mkdir -p ${BUILDDIR} || true

update-go-mod: create-build-dir
	wget -LO ${BUILDDIR}/gotify-server.mod https://raw.githubusercontent.com/gotify/server/${GOTIFY_VERSION}/go.mod
	gomod-cap -from ${BUILDDIR}/gotify-server.mod -to go.mod
	rm ${BUILDDIR}/gotify-server.mod || true
	go mod tidy

get-gotify-server-go-version: create-build-dir
	rm ${BUILDDIR}/gotify-server-go-version || true
	wget -LO ${BUILDDIR}/gotify-server-go-version https://raw.githubusercontent.com/gotify/server/${GOTIFY_VERSION}/GO_VERSION

build-linux-amd64: get-gotify-server-go-version update-go-mod
	docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:$(GO_VERSION)-linux-amd64 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-amd64${FILE_SUFFIX}.so /proj

build-linux-arm-7: get-gotify-server-go-version update-go-mod
	docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:$(GO_VERSION)-linux-arm-7 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-arm-7${FILE_SUFFIX}.so /proj

build-linux-arm64: get-gotify-server-go-version update-go-mod
	docker run --rm -v "$$PWD/.:/proj" -w /proj gotify/build:$(GO_VERSION)-linux-arm64 go build -a -installsuffix cgo -ldflags "-w -s" -buildmode=plugin -o build/${PLUGIN_NAME}-linux-arm64${FILE_SUFFIX}.so /proj

build: build-linux-arm-7 build-linux-amd64 build-linux-arm64

.PHONY: build
