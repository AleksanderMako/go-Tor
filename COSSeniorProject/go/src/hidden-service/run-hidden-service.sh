#!/bin/bash
echo $PWD
docker rm -f hiddenservice
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hiddenservice .
docker build -t hiddenservice .
docker run -it --name hiddenservice -p 5000:5000 --network=registry-app_registry-net-dev hiddenservice