function build()
{
    local exe=./service/$1
    local out=$2
    local AUTHOR=$3

    if [[ -z "$out" ]]; then
        out=bin/$1
    fi

    echo ">> 开始编译 $exe => $out"

    local BUILD_TIME=`date '+%Y-%m-%d %H:%M:%S %z'`
    local GO_VERSION=`go version`
    local prefix='duck/lazy.'
    local LDFLAGS="-X '${prefix}BUILD_TIME=${BUILD_TIME}' -X '${prefix}GO_VERSION=${GO_VERSION}' -X '${prefix}AUTHOR=${AUTHOR}' -s -w"
    
    if go build -ldflags "${LDFLAGS}" -o $out $exe; then
        echo "编译成功 $ret 1"
    else
        echo "编译失败 $ret 1"
        return 1
    fi
}

build $*