#!/bin/sh
echo "start social rpc"
cd apps/social/rpc;
goctl rpc protoc social.proto --go_out=. --go-grpc_out=. --zrpc_out=.