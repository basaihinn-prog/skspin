
RemoteUSR=zy
DATA_DIR=/home/$RemoteUSR/.rsync-build/source
GO_CACHE=/home/$RemoteUSR/.rsync-build/cache

# 目标服务器
Remote=$RemoteUSR@mega-test


# $1 编译命令
function dockerRun() 
{
    local CWD=$(basename `pwd`)
    # go docker 缓存共享目录
    local image=rsync-build

    if test -z "$(docker images | grep ${image})"; then
        echo "正在为你安装编译环境..."
        docker build -t ${image}  --build-arg USER_ID=$(id -u) --build-arg GROUP_ID=$(id -g) .
    fi

    local DIR=`realpath $(pwd)/..`
    echo $1

    docker run \
        --rm  \
        -u $(id -u):$(id -g) \
        -w /app/$CWD/ \
        -v $GO_CACHE:/gocache \
        -v $DIR:/app \
        $image \
        $1
}

function sync_and_build() {
    # 确保这些目录存在
    ssh $Remote "
        if [ ! -d $DATA_DIR ] ; then
            mkdir -p $DATA_DIR
        fi
        if [ ! -d $GO_CACHE ] ; then
            mkdir -p $GO_CACHE
        fi
        "

    # 获取当前操作人
    local GitUser=$DRONE_COMMIT_AUTHOR
    if [ ! $DRONE ]; then 
        GitUser=`git config user.name`
    fi

    echo "gen_proto"
    sh ./gen_proto.sh

    echo "同步代码"

    local CWD=$(basename `pwd`)

    if [ ! -d "bin" ] ; then
        mkdir "bin"
    fi 

    # 将需要同步的目录存入list.txt
    # 同步命令运行于当前目录的父目录
    echo "$CWD/" > bin/list.txt
    echo "duck/" >> bin/list.txt

    rsync -avczPhr --delete --max-size=1m  \
    --exclude=.git \
    --exclude=.DS_Store \
    --exclude=bin \
    --files-from=./bin/list.txt \
    ../ $Remote:$DATA_DIR/$GitUser

    echo "开始编译"
    
    local dst=$DATA_DIR/$GitUser/$CWD
    ssh $Remote "
    cd $dst
    sh rsync-build.sh _fly_ $1 $dst
    "
}

# _fly_ 特殊占位字符串

if [ $1 == "_fly_" ] ; then
    exe=$2
    dockerRun "sh ./build.sh $exe bin/$exe.temp"

    if [ ! -d ~/game/bin/ ] ; then
        mkdir -p ~/game/bin/
    fi 

    mv bin/$exe.temp ~/game/bin/
    cd ~/game/

    docker compose stop $exe
    mv ./bin/$exe.temp ./bin/$exe
    docker compose up -d $exe
    
else
    sync_and_build $1
fi