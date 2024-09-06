//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/tebeka/go2xunit"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "github.com/gogo/protobuf/gogoproto"
)