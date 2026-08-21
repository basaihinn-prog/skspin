- vim ~/caddy/Caddyfile and ./caddy reload
- scp pggateway_config.yaml {doudou-test}:/data/game/bin/config/  and modify
- vim bin/config/grpc_route.yaml
- install npm && install zx
- rsync -avc /data/game/service/fakeapp/fake-lobby/dist doudou-test:/data/game/service/fakeapp/fake-lobby/



## 购买服务器  dou-ph-prod, 添加 rsa.pub 登录
- 修改 安全组
- 	
> 目的:
> 27017 mongodb
> 50002 pggateway
> 7890 cdn source
  
rpidrgames.com  dou-idr-prod
## 购买域名 添加泛解析

## 修改 inotify 上限
$ cat /proc/sys/fs/inotify/max_user_instances
$ echo 1024 > /proc/sys/fs/inotify/max_user_instances

/etc/sysctl.conf
fs.inotify.max_user_instances = 1024


## [安装mongo](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/) 
- 
> [创建 user](https://www.mongodb.com/docs/manual/tutorial/configure-scram-client-authentication/#std-label-create-user-admin)
    // pwgen -c 64 1
    use admin
    db.createUser(
    {
        user: "myUserAdmin",
        pwd: passwordPrompt(), 
        roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" },
        { role: 'dbAdminAnyDatabase', db: 'admin' }
        ]
    }

db.grantRolesToUser(
    "myUserAdmin",
    [
      { role: "root", db: "admin" }
    ]
)



    )
    db.createUser(
    {
        user: "myUserAdmin",
        pwd: passwordPrompt(), 
        roles: [
        { role: "root", db: "admin" },
        ]
    }
    )

    - vim /etc/mongod.conf
    net:
        port: 27017
        bindIp: 0.0.0.0
    security:
        authorization: enabled
> 修改密码 db.changeUserPassword(username, password)

- db.adminCommand( { shutdown: 1 } )
- mongod -f /etc/mongod.conf --fork

# mongo rs
openssl rand -base64 756 > /root/mongo.keyfile && chmod 400  /root/mongo.keyfile
rs.initiate()
rs.add( { host: "rp-hk-dev-repl:27017", priority: 0, hidden: true, votes: 0} )
 { host: "rp-idr-prod-repl:27017", priority: 0, hidden: true, votes: 0} 
## 复制 caddy 
~/caddy/
修改解析域名

## 复制 game 
/data/game/banjia.sh (先修改host)

- 同步fake-lobby
root@dou-ph-prod:/data/game# mkdir -p /data/game/service/fakeapp/fake-lobby/
rsync -avc /data/game/service/fakeapp/fake-lobby/dist  dou-ph-prod:/data/game/service/fakeapp/fake-lobby/

- 同步 mongo  game/Games  game/Lang
- /data/pggames# zx sync-db-GamesLang.mjs  (先修改host)


- copy 自研游戏
同步 db
/data/pggames/sync-db-oldgame.mjs


同步 web asserts
root@156-241-5-141:/data/h5games# rsync -av Olympus Olympus1000 Starlight Starlight1000 StarlightChristmas  doudou-test:/data/h5games/


rsync -av Hilo doudou-prod:/data/dl/


 rsync -av 1.1.354 doudou-test:/data/h5games/shared/
 rsync -av 1.1.354 doudou-prod:/data/dl/shared/
 rsync -av 1.1.354 dou-idr-prod:/data/h5games/shared/

rsync -azvP /data/h5games/Roma  doudou-prod:/data/dl/
rsync -azvP /data/h5games/Roma  doudou-test:/data/h5games/
rsync -azvP /data/h5games/Roma  dou-idr-prod:/data/h5games/

rsync -av /data/h5games/lottery-pc/ root@doudou-test:/data/h5games/lottery-pc/


同步 icon
rsync -avc root@doudou-test:/data/dl/icon/ /data/dl/icon/

- 自研下注历史页面
root@dou-ph-prod:/data# rsync -avc doudou-prod:/data/game-history /data/
**or** exec  /data/game-history/sync.sh


## 复制PG games
- grpc_route.yaml 修改 backuphost: doudou-test 然后直接 在dev 服执行 root@156-241-5-141:/data/pggames# ./sync-restart.mjs **输入 dump 部署**


## 同步game-admin
rsync -avz -P /data/game-admin/dist root@doudou-prod:/data/game-admin/

## 同步 cache  *可选步骤*
- rsync -avc /data/game/bin/cache root@doudou-prod:/data/game/bin/
- rsync -avc root@doudou-prod:/data/game/bin/cache/m.pgsoft-games.com/custom  /data/game/bin/cache/m.pgsoft-games.com/    // copy to local 
- rsync -avc /data/game/bin/cache/m.pgsoft-games.com/custom  doudou-test:/data/game/bin/cache/m.pgsoft-games.com/  // dev copy to remote 

## 更新 apidoc/  部署不需要
cp /data/game/service/gamecenter/docs/*.md /data/apidoc/

## 修改配置  
- sync from dou-prod

> add host 47.236.107.220 doudou-prod
> root@dou-ph-prod:/data/game/bin/config# rsync --exclude='.git' -avc doudou-prod:/data/game/bin/config/ .

> pggateway_config.yaml
> jiligateway_config.yaml
> grpc_route.yaml
    - backuphost: 156.241.5.141
    - CaiPiao.http.test: 127.0.0.1:13001 (彩票需要添加)
> game_config.yaml
> ~/caddy/Caddyfile
> comm_config.yaml
> admin_setting.yaml

## 运行 /data/game/
> scp /usr/local/bin/launch dou-ph-prod:/usr/local/bin/
> YingCaiShen Olympus 需要同步db /data/pggames/sync-db-oldgame.mjs
> 
