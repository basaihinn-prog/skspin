RP Games Developer Team
## 游戏系统对接文档 Transfer Mode

游戏系统对接文档

## 接口描述
### 接口协议
- 请求数据Header设置 Content-Type: application/json
- 全部接口皆使用 POST
- 返回JSON 固定 {"code":int, error: string, data: object}  样式, code字段不为0则表示有错误发生, 此时不应再使用data字段, code:0 表示成功, code不为0时 error代表错误描述
  1. ```{"code":0,"error":"","data":{"Balance":245.9}}```
  2. ```{"code":6007,"error":"Player account does not exist","data":null}```
- datetime 用 RFC3339 格式
- ***可以使用 [mockapi](https://gamecenter.rpgamestest.com/apimock/) 测试***     
  注: 只需要关注 /api/v1/ 的接口即可

### 接口清单
- 所有请求的接口 添加2个Header
   1. AppID: 我方提供的appid 信息
   2. AppSecret: 我方提供的secret 信息
- 文档中的API URL、AppID、AppSecret将由贵司申请线路后由我方提供。

 API URL:   https://gamecenter.rpgamestest.com/
 AppID:     PartySlots
 AppSecret: 2de5c9c3-76a2-428a-aba0-XXXXXXXXXXXX

## API接口

### 创建玩家帐号(optional)
> URL: APIURL/api/v1/player/create

#### 请求参数:

| 参数名  |     类型   |         描述 |
| :----- | :--------: | -------: |
| UserID |  string[4-40]  | 运营商唯一标识 |

- eg. ```{"UserID":"abc"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| Pid |  int64  | 平台玩家唯一标识 |

- eg. ```{"error":"","data":{"Pid":100064}}```

- 重复创建同一个玩家, 返回相同的结果
- /api/v1/player/transferIn 和 /api/v1/game/launch 会自动创建玩家

### 玩家转入

> URL: APIURL/api/v1/player/transferIn

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |
| Amount |  float  | 转入金额 >0 |
| TraceId |  string[4-40]  | 交易单号 |

- eg. ```{"UserID":"user_id","Amount":123.45,"TraceId":"abc-def-gh"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| AfterBalance |  float  | 充值成功后余额 |

- eg. ```{"error":"","data":{"AfterBalance":123.45}}```
- eg. ```{"error":"Order already exists","data":null}```

### 玩家转出

> URL: APIURL/api/v1/player/transferOut

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |
| Amount |  float  | 转出金额 >0 |
| TraceId |  string[4-40]  | 交易单号 |

- eg. ```{"UserID":"user_id","Amount":1,"TraceId":"1bcfsa-dskq-req"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| AfterBalance |  float  | 充值成功后余额 |

- eg. ```{"code":0,"error":"","data":{"AfterBalance":243.9}}```
- eg. ```{"code":1,"error":"Order already exists","data":null}``` 订单号重复
- eg. ```{"code":1,"error":"Insufficient wallet balance","data":null}``` 余额不足


### 玩家转出所有余额

> URL: APIURL/api/v1/player/transferOutAll

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |
| TraceId |  string[4-40]  | 交易单号 |

- eg. ```{"UserID":"user_id","TraceId":"fdsafaagg1234"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| Amount |  float  | 共转出的总额 |

- eg. ```{"error":"","data":{"Amount":123.45}}```

### 查询订单
> URL: APIURL/api/v1/transaction/queryOrder
    
    注: 转入转出后, 如果api调用超时, 通常需要调用此接口检查订单状态


#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| TraceId |  string[4-40]  | 交易单号 |

- eg. ```{"TraceId": "fdsafaagg1234"}```

#### 返回结果:


| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| ID |  string  | 平台内部唯一ID |
| TraceId |  string  | 交易单号 |
| CreateTime |  datetime  | 订单创建时间 |
| Pid |  int64  | 平台内部玩家id |
| UserID |  string  | operator玩家id |
| AppID |  string  | 运营商标识 |
| Amount |  float  | 转入/转出总额 |
| Action |  string  | in,out,allout  转入,转出,全部转出 |
| Error |  string  | 订单执行的错误信息, 如果是"" 表示成功 |
| Completed |  bool  | 订单是否执行完成, 如果是false, 请稍等再重新使用此接口检查 |

- 正常成功的订单 ```{"error":"","data":{"Order":{"TraceId":"fdsafaagg1234","CreateTime":"2023-12-07T10:55:42.15+08:00","Pid":100065,"UserID":"user_id","AppID":"faketrans","Amount":123.45,"Action":"outall","Error":"","Completed":true}}}```
- 订单存在但是转账失败 ```{"error":"","data":{"Order":{"TraceId":"fdsafaagg1234","CreateTime":"2023-12-07T10:55:42.15+08:00","Pid":100065,"UserID":"user_id","AppID":"faketrans","Amount":123.45,"Action":"outall","Error":"some internal error occur","Completed":true}}}```
- 订单不存在 ```{"error":"Order does not exist","data":null}```


### 获取玩家余额

> URL: APIURL/api/v1/player/balance

#### 请求参数:

| 参数名  | 类型 |     描述 |
| :----- | :--: | -------: |
| UserID |  string[4-40]  | 运营商的玩家唯一标识 |

- eg. ```{"UserID":"user_id"}```

#### 返回结果:

| 参数名  | 类型 |     描述 |
| :----- | :--: | :-------: |
| Balance |  float  | 玩家当前可用余额 |
| Unsettled |  float  | 投注后尚未结算的余额 **deprecated** | 
| UnsettledDetails |  map  | 投注后尚未结算的余额详情  **deprecated** |

- 玩家共有 143.45, 投注Hilo 20 还未结算, 投注彩票80未结算, Unsettled 和 UnsettledDetails 已废弃, 请使用 APIURL/api/v1/player/pendinggames 实现有关逻辑
- eg. ```{"code":0,"error":"","data":{"Balance":43.45,"Unsettled":100,"UnsettledDetails":{"Hilo":20,"lottery":80}}}```


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
| RoundID |  string  | 回合ID, 同一局游戏该字段相同  |
| GameType |  int  | 游戏类型  0: Slots,  3: Poker,  4: Lottery  |

- eg. ```{"code":0,"error":"","data":{"Title":["ID","Pid","UserID","GameID","Bet","Win","InsertTime","AppID","Balance","WinLose","Grade","RoundID","GameType"],"List":[["66a86b88a080d69a40ad8de8",100007,"yysky1","pg_121",1,0,"2024-07-30T12:26:48.643+08:00","faketrans",201248.06,-1,51,"66a86b88a080d69a40ad8de7",0],["66a86b89a080d69a40ad8dea",100007,"yysky1","pg_121",1,1.2,"2024-07-30T12:26:49.27+08:00","faketrans",201248.26,0.2,51,"66a86b89a080d69a40ad8de9",0]],"NextStartID":"66a86b89a080d69a40ad8dea"}}```

