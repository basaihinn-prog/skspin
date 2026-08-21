#!/bin/bash
source ./header.sh
# service=$1
# if [[ -z "$service" ]]; then
#     select service in $services ; do
#         break
#     done
# fi


# sh ./gen_proto.sh

service=""

while [[ true ]]; do
    echo 'hello'
    select service in $services ; do
        break
    done

    if [[ -d ./service/$service ]]; then
        echo "env dev"
        if ! go build  -gcflags="all=-N -l"  -o bin/$service ./service/$service; then
            exit 1
        fi
    else
        echo "env prod"
    fi



    pkill -x $service
    sleep 1
    cd bin
    launch ./$service
    sleep 1
    tail  -n 100  ${service}.log
    echo "$service complete"
    cd ..
done