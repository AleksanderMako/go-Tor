#!/bin/bash

function startPeers (){
    cd onionRouting/go-torPeer
    echo $PWD
    go build ./
    ./go-torPeer
}
portBase=9000
for((i=0;i<2;i++)); do
    export PEER_PORT=$(($portBase+$i))
    (startPeers)&
    sleep 2s
    
done


