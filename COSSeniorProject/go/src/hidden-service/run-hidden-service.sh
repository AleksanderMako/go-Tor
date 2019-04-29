#!/bin/bash
echo $PWD
function startHiddenServices () {
    docker rm -f hiddenservice
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hiddenservice .
    docker build -t hiddenservice .
    docker run -d --name ${SERVICEDNS}  --env SERVICEDNS --env SERVICE_PORT --env DATA_TYPE \
    -p ${SERVICE_PORT}:${SERVICE_PORT} --network=registry-app_registry-net-dev hiddenservice
}

portBase=5000
array=( text jpeg )
for((i=0;i<2;i++)); do
    export SERVICEDNS="hiddenservice$i"
    export SERVICE_PORT=$(($portBase+$i))
    export DATA_TYPE=${array[$i]}
    echo ${DATA_TYPE}
    startHiddenServices
    sleep 5s
done
