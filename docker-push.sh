#!/usr/bin/env bash

TAG=`cat VERSION`

docker build -t mylxsw/tech-share .

docker tag mylxsw/tech-share mylxsw/tech-share:$TAG
docker tag mylxsw/tech-share:$TAG mylxsw/tech-share:latest
docker push mylxsw/tech-share:$TAG
docker push mylxsw/tech-share:latest

