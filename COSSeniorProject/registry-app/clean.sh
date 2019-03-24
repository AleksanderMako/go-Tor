#!/bin/bash


cmd=$1
if [  $cmd -eq  "0"  ]
then
    docker-compose up --build
fi

if [  $cmd -eq  "1"  ]
then
    docker rm -f $(docker ps -aq)
fi

if [  $cmd -eq  "99"  ]
then
    docker rm -f $(docker ps -aq)  || true
fi

if [  $cmd  = "t"  ]
then
    docker-compose -f docker-compose-test.yml up --build
fi