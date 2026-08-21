RP Games Developer Team
## 游戏系统对接文档 Transfer Mode

- [接口描述](#接口描述)
  - [接口协议](#接口协议)
  - [接口请求说明](#接口请求说明)
- [API接口](#api接口)
  - [创建玩家帐号(optional)](#创建玩家帐号optional)
    - [1) 请求地址](#1-请求地址)
    - [2）请求参数:](#2请求参数)
    - [3) 返回结果:](#3-返回结果)
  - [玩家转入](#玩家转入)
    - [1) 请求地址](#1-请求地址-1)
    - [2) 请求参数:](#2-请求参数)
    - [3) 返回结果:](#3-返回结果-1)
  - [玩家转出](#玩家转出)
    - [1) 请求地址](#1-请求地址-2)
    - [2) 请求参数:](#2-请求参数-1)
    - [3) 返回结果:](#3-返回结果-2)
  - [玩家转出所有余额](#玩家转出所有余额)
    - [1) 请求地址](#1-请求地址-3)
    - [2) 请求参数:](#2-请求参数-2)
    - [3) 返回结果:](#3-返回结果-3)
  - [查询订单](#查询订单)
    - [1) 请求地址](#1-请求地址-4)
    - [2) 请求参数:](#2-请求参数-3)
    - [3) 返回结果:](#3-返回结果-4)
  - [获取玩家余额](#获取玩家余额)
    - [1) 请求地址](#1-请求地址-5)
    - [2) 请求参数:](#2-请求参数-4)
    - [3) 返回结果:](#3-返回结果-5)
  - [获取玩家当前未结算的游戏列表](#获取玩家当前未结算的游戏列表)
    - [1) 请求地址](#1-请求地址-6)
    - [2) 请求参数](#2-请求参数-5)
    - [3) 返回结果](#3-返回结果-6)
  - [获取游戏列表](#获取游戏列表)
    - [1) 请求地址](#1-请求地址-7)
    - [2) 请求参数](#2-请求参数-6)
    - [3) 返回结果](#3-返回结果-7)
  - [获取游戏登录URL](#获取游戏登录url)
    - [1) 请求地址](#1-请求地址-8)
    - [2) 请求参数:](#2-请求参数-7)
    - [3) 返回结果:](#3-返回结果-8)
  - [拉取下注历史](#拉取下注历史)
    - [1) 请求地址](#1-请求地址-9)
    - [2) 请求参数](#2-请求参数-8)
    - [3) 返回结果](#3-返回结果-9)
  - [语言和标识符对应关系](#语言和标识符对应关系)

## 接口描述
### 接口协议
- 请求数据Header设置 Content-Type: application/json
- 全部接口皆使用 POST
- 返回JSON 固定 {"code":int, error: string, data: object}  样式, code字段不为0则表示有错误发生, 此时不应再使用data字段, code:0 表示成功, code不为0时 error代表错误描述
  1. ```{"code":0,"error":"","data":{"Balance":245.9}}```
  2. ```{"code":6007,"error":"Player account does not exist","data":null}```
  3. ```{"code":6008,"error":"Player status is close","data":null}```
  4. ```{"code":6010,"error":"Access denied","data":null}```
- datetime 用 RFC3339 格式
- ***可以使用 [mockapi](https://gamecenter.rpgamestest.com/apimock/) 测试***     
  注: 只需要关注 /api/v1/ 的接口即可

### 接口请求说明
- 所有请求的接口 添加2个Header
  1. AppID: 我方提供的appid 信息
  2. AppSecret: 我方提供的secret 信息
- 文档中的API URL、AppID、AppSecret将由贵司申请线路后由我方提供。

API URL:   https://gamecenter.rpgamestest.com/
AppID:     PartySlots
AppSecret: 2de5c9c3-76a2-428a-aba0-XXXXXXXXXXXX

## API接口

### 创建玩家帐号(optional)

#### 1) 请求地址
> URL: APIURL/api/v1/player/create

#### 2）请求参数:
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "abc"
}
```

#### 3) 返回结果:
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Pid</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int64</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">平台玩家唯一标识</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "Pid": 100064
  }
}
```

- 重复创建同一个玩家, 返回相同的结果
- /api/v1/player/transferIn 和 /api/v1/game/launch 会自动创建玩家

### 玩家转入

#### 1) 请求地址
> URL: APIURL/api/v1/player/transferIn

#### 2) 请求参数:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Amount</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">转入金额 >0 </td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">TraceId</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">交易单号</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "user_id",
  "Amount": 123.45,
  "TraceId": "abc-def-gh"
}
```

#### 3) 返回结果:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">AfterBalance</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">充值成功后余额</td>
  </tr>
</table>

- 正常返回
- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "AfterBalance": 123.45
  }
}
```
- 订单重复
- <font color=#FF000 >示例</font>：
```json
{
  "error": "Order already exists",
  "data": null
}
```

### 玩家转出

#### 1) 请求地址
> URL: APIURL/api/v1/player/transferOut

#### 2) 请求参数:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Amount</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">转入金额 >0 </td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">TraceId</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">交易单号</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "user_id",
  "Amount": 1,
  "TraceId": "1bcfsa-dskq-req"
}
```

#### 3) 返回结果:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">AfterBalance</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">转出成功后余额</td>
  </tr>
</table>

- 正常返回:
- <font color=#FF000 >示例</font>：
```json
{
  "code": 0,
  "error": "",
  "data": {
    "AfterBalance": 243.9
  }
}
```

- 订单号重复
- <font color=#FF000 >示例</font>：
```json
{
  "code": 1,
  "error": "Order already exists",
  "data": null
}
```

- 余额不足
- <font color=#FF000 >示例</font>：
```json
{
  "code": 1,
  "error": "Insufficient wallet balance",
  "data": null
}
```


### 玩家转出所有余额

#### 1) 请求地址
> URL: APIURL/api/v1/player/transferOutAll

#### 2) 请求参数:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">TraceId</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">交易单号</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "user_id",
  "TraceId": "fdsafaagg1234"
}
```

#### 3) 返回结果:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Amount</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">共转出的总额</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "Amount": 123.45
  }
}
```

### 查询订单

#### 1) 请求地址
> URL: APIURL/api/v1/transaction/queryOrder

    注: 转入转出后, 如果api调用超时, 通常需要调用此接口检查订单状态


#### 2) 请求参数:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">TraceId</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">交易单号</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "TraceId": "fdsafaagg1234"
}
```

#### 3) 返回结果:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">ID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">平台内部唯一ID</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">TraceId</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">交易单号</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">CreateTime</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">datetime</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">订单创建时间</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Pid</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int64</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">平台内部玩家id</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">operator玩家id</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">AppID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商标识</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Amount</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">转入/转出总额</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Action</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">in,out,allout  转入,转出,全部转出</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Error</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">订单执行的错误信息, 如果是"" 表示成功</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Completed</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">bool</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">订单是否执行完成, 如果是false, 请稍等再重新使用此接口检查</td>
  </tr>
</table>

- 正常成功的订单
- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "Order": {
      "TraceId": "fdsafaagg1234",
      "CreateTime": "2023-12-07T10:55:42.15+08:00",
      "Pid": 100065,
      "UserID": "user_id",
      "AppID": "faketrans",
      "Amount": 123.45,
      "Action": "outall",
      "Error": "",
      "Completed": true
    }
  }
}
```
- 订单存在但是转账失败
- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "Order": {
      "TraceId": "fdsafaagg1234",
      "CreateTime": "2023-12-07T10:55:42.15+08:00",
      "Pid": 100065,
      "UserID": "user_id",
      "AppID": "faketrans",
      "Amount": 123.45,
      "Action": "outall",
      "Error": "some internal error occur",
      "Completed": true
    }
  }
}
```
- 订单不存在
- <font color=#FF000 >示例</font>：
```json
{
  "error": "Order does not exist",
  "data": null
}
```


### 获取玩家余额

#### 1) 请求地址
> URL: APIURL/api/v1/player/balance

#### 2) 请求参数:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
</table>

- eg. ```{"UserID":"user_id"}```

#### 3) 返回结果:

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Balance</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">玩家当前可用余额</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Unsettled</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">投注后尚未结算的余额 **deprecated**</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UnsettledDetails</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">map</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">投注后尚未结算的余额详情  **deprecated**</td>
  </tr>
</table>

- 玩家共有 143.45, 投注Hilo 20 还未结算, 投注彩票80未结算, Unsettled 和 UnsettledDetails 已废弃, 请使用 APIURL/api/v1/player/pendinggames 实现有关逻辑
- <font color=#FF000 >示例</font>：
```json
{
  "code": 0,
  "error": "",
  "data": {
    "Balance": 43.45,
    "Unsettled": 100,
    "UnsettledDetails": {
      "Hilo": 20,
      "lottery": 80
    }
  }
}
```


### 获取玩家当前未结算的游戏列表

#### 1) 请求地址
> URL: APIURL/api/v1/player/pendinggames

#### 2) 请求参数
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "user_id"
}
```

#### 3) 返回结果

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Balance</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">玩家当前可用余额</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">PendingGames</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">投注后尚未结算的游戏列表</td>
  </tr>
</table>

- [Balance]只对转账模式有意义, [PendingGames] 返回下注后还未完成全部结算的游戏列表
- <font color=#FF000 >示例</font>：
```json
{
  "code": 0,
  "error": "",
  "data": {
    "Balance": 4259.94,
    "PendingGames": [
      "lottery",
      "pg_1489936",
      "pg_20",
      "pg_37",
      "pg_58",
      "pg_59",
      "pg_94"
    ]
  }
}
```

### 获取游戏列表

#### 1) 请求地址
> URL: APIURL/api/v1/game/list

#### 2) 请求参数

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Language</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">语言设定 th, en ...</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>:
```json
{
    "Language": "th"
}
```

#### 3) 返回结果

<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">ID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏ID</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Name</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏名字</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Type</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">0: 拉霸游戏, 1: 捕鱼游戏, 3: 棋牌游戏, 4: 彩票游戏</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">IconUrl</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏icon</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "List": [
      {
        "ID": "XingYunXiang",
        "Name": "Ganesha Fortune",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/XinYunXiang.png"
      },
      {
        "ID": "YingCaiShen",
        "Name": "Caishen Wins",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/YingCaiShen.png"
      },
      {
        "ID": "NiuBi",
        "Name": "NiuBi",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/NiuBi.png"
      },
      {
        "ID": "BaoZang",
        "Name": "Treasures of Aztec",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/BaoZang.png"
      },
      {
        "ID": "ZhaoCaiMao",
        "Name": "Lucky  Neko",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/ZhaoCaiMao.png"
      },
      {
        "ID": "Roma",
        "Name": "Roma",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/Roma.png"
      },
      {
        "ID": "RomaX",
        "Name": "RomaX",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/RomaX.png"
      },
      {
        "ID": "TuZi",
        "Name": "Fortune Rabbit",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/JinQianTu.png"
      },
      {
        "ID": "JinNiu",
        "Name": "Fortune OX",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/JinNiu.png"
      },
      {
        "ID": "MaJiang",
        "Name": "Mahjong Ways",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/MaJiang.png"
      },
      {
        "ID": "MaJiang2",
        "Name": "Mahjong Ways2",
        "Type": 0,
        "IconUrl": "https://dl.kafa010.com/icon/MaJiang2.png"
      },
      {
        "ID": "Hilo",
        "Name": "Hilo",
        "Type": 3,
        "IconUrl": ""
      }
    ]
  }
}
```



### 获取游戏登录URL

#### 1) 请求地址
> URL: APIURL/api/v1/game/launch

#### 2) 请求参数:
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
   <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">GameID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏ID</td>
  </tr>
   <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Language</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">语言设定 th, en ...</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "UserID": "operator_user_abcd",
  "GameID": "XingYunXiang",
  "Language": "th"
}
```

#### 3) 返回结果:
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Url</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏登录URL</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "error": "",
  "data": {
    "Url": "https://h5games.rpgamestest.com/XingYunXiang/index.html?l=th&t=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxODcsIlMiOjEwMDIsIkQiOiJYaW5nWXVuWGlhbmcifQ.9td-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxmw"
  }
}
```


### 拉取下注历史

#### 1) 请求地址
> URL: APIURL/api/v1/record/betlist

#### 2) 请求参数
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">StartID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[24]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">起始的下注记录ID, 结果不会包含此ID记录</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Count</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">拉取数量, 范围 1~5000, 超出将会限制到 1或5000</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "StartID": "656fd4993be8edaa4d37830e",
  "Count": 2
}
```

连续拉取的时候, 下一次拉起请传入上次拉取的最后一个记录的ID

#### 3) 返回结果
<table style="border-collapse: collapse; width: 100%;">
  <tr style="background-color: #000000; color: #ffffff;">
    <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">参数名</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">类型</th>
    <th style="border: 1px solid #ddd; padding: 8px; text-align: center;">描述</th>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">ID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">记录的唯一ID 升序排列</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Pid</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int64</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">平台玩家ID</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">UserID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string[4-40]</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商的玩家唯一标识</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">GameID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">游戏ID</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Bet</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">下注总额</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Win</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">赢分总额</td>
  </tr> 
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">WinLose</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">净输赢</td>
  </tr> 
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">InsertTime</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">datetime</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">记录入库的时间</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">AppID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">运营商ID</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Balance</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">float</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">结算后余额</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">Grade</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">int</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">投注挡位</td>
  </tr>
  <tr style="background-color:#f2f2f2;">
    <td style="background-color: #000000; color: #ffffff; border: 1px solid #ddd; padding: 8px; text-align: left;">RoundID</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">string</td>
    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">回合ID, 同一局游戏该字段相同</td>
  </tr>
</table>

- <font color=#FF000 >示例</font>：
```json
{
  "code": 0,
  "error": "",
  "data": {
    "Title": [
      "ID",
      "Pid",
      "UserID",
      "GameID",
      "Bet",
      "Win",
      "InsertTime",
      "AppID",
      "Balance",
      "WinLose",
      "Grade"
    ],
    "List": [
      [
        "65dda8a2877e2ea628b7a5aa",
        100034,
        "operator_user_abcd",
        "Hilo",
        5,
        10,
        "2024-02-27T16:17:22.429+07:00",
        "faketrans",
        128.45,
        5,
        -1
      ],
      [
        "65dda8ca877e2ea628b7a5ab",
        100034,
        "operator_user_abcd",
        "Hilo",
        20,
        60,
        "2024-02-27T16:18:02.429+07:00",
        "faketrans",
        168.45,
        40,
        -1
      ]
    ]
  }
}
```

### 语言和标识符对应关系

- en 英文
- da 丹麦文
- de 德文
- es 西班牙文
- fi 芬兰文
- fr 法文
- id 印尼文
- it 意大利文
- ja 日文
- ko 韩文
- nl 荷兰文
- no 挪威文
- pl 波兰文
- pt 葡萄牙文
- ro 罗马尼亚文
- ru 俄文
- sv 瑞典文
- th 泰文
- tr 土耳其文
- vi 越南文
- my 缅甸文
