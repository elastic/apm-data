#!/bin/sh

TOOLS_DIR=$(dirname "$(readlink -f -- "$0")")

PATH="${TOOLS_DIR}/build/bin:${PATH}" protoc --proto_path=./model/proto/ --go_out=. --go_opt=module=github.com/elastic/apm-data --go-vtproto_out=. --go-vtproto_opt=features=marshal+unmarshal+size+pool,module=github.com/elastic/apm-data ./model/proto/*.proto
