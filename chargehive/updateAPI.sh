#!/usr/bin/env bash
echo "Downloading api.proto from chargehive:"
curl -O https://api.chargehive.com/api.proto

echo
echo "Creating Go api from proto.api:"
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/chargehive/proto \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  --gogo_out=plugins=grpc:.\
  api.proto

echo "Updated!"