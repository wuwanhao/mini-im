#!/bin/sh
echo "start social api"
cd ../apps/social/api;
goctl api go -api social.api -dir . -style gozero
