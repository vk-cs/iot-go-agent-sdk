CWD=$(shell pwd)
USER=$(shell id -u)
GROUP=$(shell id -g)
GOPATH=$(shell go env GOPATH)
HTTP_CLIENT_PATH=./gen/swagger/http_client

.PHONY: generate
generate:
	mkdir -p ${HTTP_CLIENT_PATH} || true &&\
		rm -rf ${HTTP_CLIENT_PATH}/* || true && \
		docker run \
			--rm \
			-it \
			--user ${USER}:${GROUP} \
			-e GOPATH=${GOPATH}:/go \
			-v ${HOME}:${HOME} \
			-w ${CWD} \
			quay.io/goswagger/swagger \
			generate client -f ./swagger.json --target=${HTTP_CLIENT_PATH} -A http

.PHONY: prepare-bin
prepare-bin:
	rm -rf ./bin || true
	mkdir -p ./bin || true

.PHONY: build-linux-example
build-linux-example:
	GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/mqtt_example ./cmd/mqtt_example

.PHONY: build-darwin-example
build-darwin-example:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/mqtt_example ./cmd/mqtt_example

.PHONY: build-windows-example
build-windows-example:
	GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64/mqtt_example.exe ./cmd/mqtt_example

.PHONY: build-example
build-example: \
	prepare-bin \
	build-linux-example \
	build-darwin-example \
	build-windows-example
