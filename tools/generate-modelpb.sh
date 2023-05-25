#!/bin/sh

TOOLS_DIR=$(dirname "$(readlink -f -- "$0")")

PATH="${PATH}:${TOOLS_DIR}/bin" protoc --proto_path=./model/proto/ --go_out=. --go_opt=module=github.com/elastic/apm-data --go-vtproto_out=. --go-vtproto_opt=features=marshal+unmarshal+size,module=github.com/elastic/apm-data ./model/proto/*.proto
