#!/bin/bash

host=""
hosts="
127.0.0.1
"
select host in $hosts ; do
    break
done

echo host $host

source ./header.sh
bin=""

while [[ true ]]; do
    select bin in $services ; do
        break
    done

    echo building... $bin

    # if ! go build -gcflags="all=-N -l"  -o bin/ ./service/$bin; then
    if ! go build -o bin/ ./service/$bin; then
        exit 1
    fi

    rsync -vcPz bin/${bin} root@${host}:/data/game/bin/${bin}.new \
    && \
    ssh root@${host} "
    cd /data/game/bin/
    pkill -x ${bin}
    sleep 1
    mv ${bin} ${bin}.bak
    mv ${bin}.new ${bin}
    launch ./${bin}
    sleep 1
    tail -f ${bin}.log
    "
    echo "$bin complete"
    # setsid ./${bin} &> ${bin}.log &

done