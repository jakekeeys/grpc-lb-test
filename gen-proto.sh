#!/usr/bin/env sh

# GRPC
protoc -I sample/ \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=plugins=grpc:sample \
    sample/sample.proto

# Gateway
protoc -I sample/ \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:sample \
    sample/sample.proto

# Swagger
protoc -I sample/ \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:sample \
    sample/sample.proto
cp sample/sample.swagger.json swagger.json