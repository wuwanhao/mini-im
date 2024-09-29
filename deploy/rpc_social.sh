#!/bin/sh
echo "start user rpc"
cd apps/user/rpc;
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.