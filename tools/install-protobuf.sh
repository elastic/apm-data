#!/bin/sh

PLATFORM=$(go env GOOS)
PROTOBUF_VERSION="23.1"
PROTOC_GO_VERSION="v1.30.0"
VTPROTOBUF_VERSION="v0.4.0"


if [ "${PLATFORM}" = "darwin" ]; then
	PROTOBUF_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-osx-x86_64.zip"
elif [ "${PLATFORM}" = "linux" ]; then 
	PROTOBUF_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip"
elif [ "${PLATFORM}" = "windows" ]; then 
	PROTOBUF_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-win64.zip"
else
	echo "Unsupported platform: ${PLATFORM}"
	exit 1
fi

TOOLS_DIR=$(dirname "$(readlink -f -- "$0")")
PROTOBUF_ZIP="/tmp/protobuf.zip"

curl -L "${PROTOBUF_URL}" -o "${PROTOBUF_ZIP}"

if ! unzip -j -o "${PROTOBUF_ZIP}" "bin/protoc" -d "${TOOLS_DIR}/bin"; then
	echo "failed to extract protobuf"
	exit 1
fi

if ! PATH="${TOOLS_DIR}/bin" protoc --version; then 
	echo "failed to verify protobuf installation"
	exit 1
fi

GOBIN="${TOOLS_DIR}/bin" go install "google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GO_VERSION}"
GOBIN="${TOOLS_DIR}/bin" go install "github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@${VTPROTOBUF_VERSION}"
