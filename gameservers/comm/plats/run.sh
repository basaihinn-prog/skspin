#!/bin/bash
service=plats
if ! go build  -gcflags="all=-N -l"  -o ../../bin/; then
    exit 1
fi


cd ../../bin

pkill -x $service
sleep 1
# setsid ./$service &> ${service}.log &
# /data/game/tool/launch/launch ./$service
launch ./$service
sleep 1
tail -f -n 100  ${service}.log