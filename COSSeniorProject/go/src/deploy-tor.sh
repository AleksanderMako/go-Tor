#!/bin/bash
# the function compiles go libs statically
# this way the binaries may be run inside a scratch contianer

function startPeers (){
    cd onionRouting/go-torPeer
    echo $PWD
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
    docker build -t peer .
    docker  run -it --env DNS --env PEER_PORT  -p ${PEER_PORT}:${PEER_PORT}   --network=registry-app_registry-net-dev peer
}
# portBase=9000
# for((i=0;i<2;i++)); do
#     export PEER_PORT=$(($portBase+$i))
#     (startPeers)&
#     sleep 2s

# done

export PEER_PORT=5500
export DNS=peer
(startPeers)

