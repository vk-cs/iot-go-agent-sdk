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
