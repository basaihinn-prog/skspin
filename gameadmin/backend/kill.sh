#!/bin/bash
source ./header.sh
service=$1
if [[ -z "$service" ]]; then
    select service in $services ; do
        break
    done
fi

pkill -x $service