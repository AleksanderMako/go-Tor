#!/bin/bash
# the function compiles go libs statically
# this way the binaries may be run inside a scratch contianer

function startPeers (){
    cd onionRouting/go-torPeer
    echo $PWD
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
    docker build -t peer .
    docker  run -d --env DNS --env PEER_PORT  -p ${PEER_PORT}:${PEER_PORT}  --name $DNS --network=registry-app_registry-net-dev peer
}

function startClient() {
    cd onionRouting/go-torClient
    echo $PWD
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client .
    docker build -t client .
    docker run -it --name torclient -p 8000:8000 --network=registry-app_registry-net-dev client
}
portBase=9000
for((i=0;i<4;i++)); do
    export PEER_PORT=$(($portBase+$i))
    export DNS="peer$i"
    (startPeers)
    sleep 2s
    
done
sleep 5s
(startClient)


# export PEER_PORT=5500

# (startPeers)

