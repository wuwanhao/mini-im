#!/bin/sh
echo "start user api"
cd apps/user/api;
goctl api go -api user.api -dir . -style gozero
