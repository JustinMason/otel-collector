#!/usr/bin/env bash
#
# Generate all protobuf bindings.
# Run from repository root.
set -e
set -u

if ! [[ "$0" =~ "scripts/genproto.sh" ]]; then
	echo "must be run from repository root"
	exit 255
fi

# if ! [[ $(protoc --version) =~ "3.15.8" ]]; then
# 	echo "could not find protoc 3.15.8, is it installed + in PATH?"
# 	exit 255
# fi

# Since we run go install, go mod download, the go.sum will change.
# Make a backup.
cp go.sum go.sum.bak

INSTALL_PKGS="golang.org/x/tools/cmd/goimports github.com/gogo/protobuf/protoc-gen-gogofast github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
for pkg in ${INSTALL_PKGS}; do
    GO111MODULE=on go install "$pkg"
done

PROM_ROOT="${PWD}"
PROM_PATH="${PROM_ROOT}/receiver/jsonlogreceiver/proto"
GOGOPROTO_ROOT="$(GO111MODULE=on go list -mod=readonly -f '{{ .Dir }}' -m github.com/gogo/protobuf)"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}:${GOGOPROTO_ROOT}/protobuf"
GRPC_GATEWAY_ROOT="$(GO111MODULE=on go list -mod=readonly -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)"

DIRS=${PROM_PATH}

echo "generating code"
for dir in ${DIRS}; do
	pushd ${dir}
		#echo ${dir}
		protoc --gogofast_out=plugins=grpc:. -I=. \
            -I="${GOGOPROTO_PATH}" \
            -I="${PROM_PATH}" \
            -I="${GRPC_GATEWAY_ROOT}/third_party/googleapis" \
            ./*.proto
		# google/protobuf is github.com/gogo/protobuf/types
		# protoc --gogofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,paths=source_relative:. -I=. \
        #     -I="${GOGOPROTO_PATH}" \
        #     ./io/prometheus/client/*.proto
		#sed -i.bak -E 's/import _ \"github.com\/gogo\/protobuf\/gogoproto\"//g' -- *.pb.go
		#sed -i.bak -E 's/import _ \"google\/protobuf\"//g' -- *.pb.go
		#sed -i.bak -E 's/\t_ \"google\/protobuf\"//g' -- *.pb.go
		#sed -i.bak -E 's/golang\/protobuf\/descriptor/gogo\/protobuf\/protoc-gen-gogo\/descriptor/g' -- *.go
		#sed -i.bak -E 's/golang\/protobuf/gogo\/protobuf/g' -- *.go
		rm -f -- *.bak
		# goimports -w ./*.go ./io/prometheus/client/*.go
	popd
done

mv go.sum.bak go.sum
