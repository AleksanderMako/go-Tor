#!/bin/bash
echo $PWD
docker rm -f torclient
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client .
docker build -t client .
docker run -it --name torclient -p 8000:8000 --network=registry-app_registry-net-dev client