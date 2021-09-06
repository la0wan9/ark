//go:build tool

package tool

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/cosmtrek/air"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "github.com/tomwright/dasel/cmd/dasel"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
