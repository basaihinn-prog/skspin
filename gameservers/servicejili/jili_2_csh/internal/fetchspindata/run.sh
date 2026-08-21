#!/bin/bash
service=fetch.out

if ! go build -o $service; then
    exit 1
fi

setsid ./$service &> fetch.log &
sleep 1
tail -f -n 100  fetch.log
