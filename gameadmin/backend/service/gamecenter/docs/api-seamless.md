RP Games Developer Team
## 游戏系统对接文档 Seamless Mode

游戏系统对接文档

##  RP平台提供的API
###  接口协议
- 请求数据Header设置 Content-Type: application/json
- 全部接口皆使用 POST
- 返回JSON 固定 {"code":int, error: string, data: object}  样式, code字段不为0则表示有错误发生, 此时不应再使用data字段, code:0 表示成功, code不为0时 error代表错误描述
  1. ```{"code":0,"error":"","data":{"Balance":245.9}}```
  2. ```{"code":6007,"error":"Player account does not exist","data":null}```
- datetime 用 RFC3339 格式
- ***可以使用 [mockapi](https://gamecenter.rpgamestest.com/apimock/) 测试***     
  注: 只需要关注 /api/v1/ 的接口即可

###  接口清单
- 所有请求的接口 添加2个Header
   1. AppID: 我方提供的appid 信息
   2. AppSecret: 我方提供的secret 信息
- 文档中的API URL、AppID、AppSecret将由贵司申请线路后由我方提供。

 API URL:   https://gamecenter.rpgamestest.com/
 AppID:     PartySlots
 AppSecret: 2de5c9c3-76a2-428a-aba0-XXXXXXXXXXXX

##  API接口

###  创建玩家帐号(optional)
> URL: APIURL/api/v1/player/create

#### 请求参数:

| 参数名  |     类型   |         描述 |
| :----- | :--------: | -------: |
| UserID |  string[4-40]  | 运营商唯一标识 |

- 可选的接口, 当玩家第一次启动游戏时, 会自动创建

- eg. ```{"UserID":"abc"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| Pid |  int64  | 平台玩家唯一标识 |

- eg. ```{"error":"","data":{"Pid":100064}}```

- 重复创建同一个玩家, 返回相同的结果
- /api/v1/player/transferIn 和 /api/v1/game/launch 会自动创建玩家

### 获取玩家当前未结算的游戏列表

> URL: APIURL/api/v1/player/pendinggames

| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |

- eg. ```{"UserID":"user_id"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| Balance |  float  | 玩家当前可用余额 |
| PendingGames |  string[]  | 投注后尚未结算的游戏列表 | 

- [Balance]只对转账模式有意义, [PendingGames] 返回下注后还未完成全部结算的游戏列表
- eg. ```{"code":0,"error":"","data":{"Balance":4259.94,"PendingGames":["lottery","pg_1489936","pg_20","pg_37","pg_58","pg_59","pg_94"]}}```

### 设置玩家RTP

> URL: APIURL/api/v1/player/setRtp

| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |
| Rtp |  int  | [0,200] rtp回报率,特殊值-1取消设置 |

- eg. ```{"UserID":"user_id","Rtp":120}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| Pid |  int64  | 平台玩家ID |
| Rtp |  int  | 上一次的rtp回报率 | 

- eg. ```{"code":0,"error":"","data":{"Pid":100207,"Rtp":50}}```



### 获取游戏列表

> URL: APIURL/api/v1/game/list

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |

- eg. ```{}``` //传入空的json 对象

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| ID |  string  | 游戏ID |
| Name |  string  | 游戏名字 |
| Type |  int    | 0: 拉霸游戏, 1: 捕鱼游戏, 3: 棋牌游戏, 4: 彩票游戏 |
| IconUrl |  string  |  游戏icon |

- eg. ```{"error":"","data":{"List":[{"ID":"XingYunXiang","Name":"Ganesha Fortune","Type":0,"IconUrl":"https://dl.kafa010.com/icon/XinYunXiang.png"},{"ID":"YingCaiShen","Name":"Caishen Wins","Type":0,"IconUrl":"https://dl.kafa010.com/icon/YingCaiShen.png"},{"ID":"NiuBi","Name":"NiuBi","Type":0,"IconUrl":"https://dl.kafa010.com/icon/NiuBi.png"},{"ID":"BaoZang","Name":"Treasures of Aztec","Type":0,"IconUrl":"https://dl.kafa010.com/icon/BaoZang.png"},{"ID":"ZhaoCaiMao","Name":"Lucky  Neko","Type":0,"IconUrl":"https://dl.kafa010.com/icon/ZhaoCaiMao.png"},{"ID":"Roma","Name":"Roma","Type":0,"IconUrl":"https://dl.kafa010.com/icon/Roma.png"},{"ID":"RomaX","Name":"RomaX","Type":0,"IconUrl":"https://dl.kafa010.com/icon/RomaX.png"},{"ID":"TuZi","Name":"Fortune Rabbit","Type":0,"IconUrl":"https://dl.kafa010.com/icon/JinQianTu.png"},{"ID":"JinNiu","Name":"Fortune OX","Type":0,"IconUrl":"https://dl.kafa010.com/icon/JinNiu.png"},{"ID":"MaJiang","Name":"Mahjong Ways","Type":0,"IconUrl":"https://dl.kafa010.com/icon/MaJiang.png"},{"ID":"MaJiang2","Name":"Mahjong Ways2","Type":0,"IconUrl":"https://dl.kafa010.com/icon/MaJiang2.png"},{"ID":"Hilo","Name":"Hilo","Type":3,"IconUrl":""}]}}```



### 获取游戏登录URL

> URL: APIURL/api/v1/game/launch

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |
| GameID |  string  | 游戏ID |
| Language |  string  | 语言设定 th, en ... |

> en 英文
> da 丹麦文
> de 德文
> es 西班牙文
> fi 芬兰文
> fr 法文
> id 印尼文
> it 意大利文
> ja 日文
> ko 韩文
> nl 荷兰文
> no 挪威文
> pl 波兰文
> pt 葡萄牙文
> ro 罗马尼亚文
> ru 俄文
> sv 瑞典文
> th 泰文
> tr 土耳其文
> vi 越南文
> my 缅甸文



- eg. ```{"UserID":"operator_user_abcd","GameID":"XingYunXiang","Language":"th"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| Url |  string  | 游戏登录URL |

- eg. ```{"error":"","data":{"Url":"https://h5games.rpgamestest.com/XingYunXiang/index.html?l=th&t=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxODcsIlMiOjEwMDIsIkQiOiJYaW5nWXVuWGlhbmcifQ.9td-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxmw"}}```


### 拉取下注历史

> URL: APIURL/api/v1/record/betlist

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| StartID |  string[24]  | 起始的下注记录ID, 结果不会包含此ID记录 |
| Count |  int  | 拉取数量, 范围 1~5000, 超出将会限制到 1或5000 |

- eg. ```{"StartID":"656fd4993be8edaa4d37830e","Count":2}```

  连续拉取的时候, 下一次拉起请传入上次拉取的最后一个记录的ID

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| ID |  string  | 记录的唯一ID 升序排列 |
| Pid |  int64  | 平台玩家ID |
| UserID |  string  | 运营商玩家ID |
| GameID |  string  |  游戏ID |
| Bet |  float  | 下注总额 |
| Win |  float  | 赢分总额 |
| WinLose |  float  | 净输赢 |
| InsertTime |  datetime  | 记录入库的时间 |
| AppID |  string  | 运营商ID |
| Balance |  float  | 结算后余额 |
| Grade |  int  | 投注挡位 |

- eg.  ```
{"code":0,"error":"","data":{"Title":["ID","Pid","UserID","GameID","Bet","Win","InsertTime","AppID","Balance","WinLose","Grade"],"List":[["65dda8a2877e2ea628b7a5aa",100034,"operator_user_abcd","Hilo",5,10,"2024-02-27T16:17:22.429+07:00","faketrans",128.45,5,-1],["65dda8ca877e2ea628b7a5ab",100034,"operator_user_abcd","Hilo",20,60,"2024-02-27T16:18:02.429+07:00","faketrans",168.45,40,-1]]}} ```



## 运营商提供的API
- OPURL 是你方提供的http访问地址, 我方会通过 http POST 请求玩家数据
- 所有请求和响应均使用Json编码
- 请首先检查AppID, 和AppSecret字段和你方匹配
- 如果你方有error发生, 请将 *code* 设置为非0值, 同时 *error* 设置为非空字符串
- eg. OPURL=https://balabala.com/xxxx or http://10.0.0.1:8888 

### 获取玩家余额

> URL: OPURL/Cash/Get

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| AppID |  string  | 运营商唯一标识 |
| AppSecret |  string  | 运营商AppSecret |
| UserID |  string  | 运营商的玩家唯一标识 |

- eg. ```{"UserID":"testuser1","AppID":"fake","AppSecret":"1234-abcd-xxxkkk"}```

#### 返回结果:
| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| Balance |  float  | 玩家当前可用余额 |

- eg. ```{"code":0,"error":"","data":{"Balance":10994.12}}```


### 修改玩家余额

> URL: OPURL/Cash/TransferInOut

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| AppID |  string  | 运营商唯一标识 |
| AppSecret |  string  | 运营商AppSecret |
| UserID |  string  | 运营商的玩家唯一标识 |
| TransactionID |  string  | 交易订单号 |
| Amount |  float  | 增加/扣除金额 (+ 增加, - 扣除) |


- 如果我方在请求此接口时发生Timeout错误, 可能会在接下来的一小段时间内连续请求若干次, 请使用 TransactionID 判断重复, 保证
- eg. ```{"UserID":"testuser1","AppID":"fake","AppSecret":"1234-abcd-xxxkkk","TransactionID":"abc-xxx-yyyy-zzz","Amount":10}```


#### 返回结果:
| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| Balance |  float  | 修改后玩家当前可用余额 |


- eg. ```{"code":0,"error":"","data":{"Balance":11004}}```
- eg. ```{"code":1,"error":"some error occur!!"}```
