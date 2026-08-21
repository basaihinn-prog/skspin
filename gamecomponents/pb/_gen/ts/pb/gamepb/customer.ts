/* eslint-disable */
import type { Empty } from "../empty";


export enum GameType {
  /** Slot - 拉霸游戏 */
  Slot = 0,
  /** Fish - 捕鱼游戏 */
  Fish = 1,
  /** Poker - 棋牌游戏 */
  Poker = 3,
}

export enum GameStatus {
  /** Open - 正常 */
  Open = 0,
  /** Maintenance - 维护中（列表中出现） */
  Maintenance = 1,
  /** Hide - 隐藏（列表中不出现） */
  Hide = 2,
}

export interface DocBetLog {
  /** @gotags:bson:"_id" */
  ID?:
    | string
    | undefined;
  /** 内部玩家ID */
  Pid?:
    | number
    | undefined;
  /** 外部玩家ID */
  UserID?:
    | string
    | undefined;
  /** 游戏ID */
  GameID?:
    | string
    | undefined;
  /** 对局ID，方便查询 */
  RoundID?:
    | string
    | undefined;
  /** 下注 */
  Bet?:
    | number
    | undefined;
  /** 输赢 */
  Win?:
    | number
    | undefined;
  /** 注释 */
  Comment?:
    | string
    | undefined;
  /** 数据插入时间 */
  InsertTime?:
    | string
    | undefined;
  /**  */
  AppID?:
    | string
    | undefined;
  /** 余额 */
  Balance?: number | undefined;
}

export interface DocBetLogList {
  List?: DocBetLog[] | undefined;
}

export interface GetLogReq {
  From?: string | undefined;
}

export interface LoginReq {
  /** 玩家id */
  UserID?:
    | string
    | undefined;
  /** 游戏ID */
  GameID?:
    | string
    | undefined;
  /** 游戏语言 */
  Language?: string | undefined;
}

export interface LoginResp {
  /** 游戏启动链接 */
  Url?: string | undefined;
}

export interface Game {
  /** @gotags:bson:"_id" */
  ID?: string | undefined;
  Name?: string | undefined;
  Type?: GameType | undefined;
  Status?: GameStatus | undefined;
}

export interface GameListResp {
  List?: Game[] | undefined;
}

export interface ModifyGoldUidReq {
  /** 玩家ID */
  UserID?:
    | string
    | undefined;
  /** 金币变化：>0 加钱， <0 扣钱 */
  Change?:
    | number
    | undefined;
  /** 注释 */
  Comment?: string | undefined;
}

export interface GetBalanceUidReq {
  UserID?: string | undefined;
}

export interface GetBalanceUidResp {
  Balance?: number | undefined;
}

/** 给接入者调用 */
export class ForCustomerRpc {
  /** 获取游戏列表 */
  static async GameList(client, req : Empty) : Promise<[GameListResp,any]> {
    return await client.send("ForCustomerRpc/GameList", req)
  }
  /** 获取游戏启链接 */
  static async Login(client, req : LoginReq) : Promise<[LoginResp,any]> {
    return await client.send("ForCustomerRpc/Login", req)
  }
  /**
   * 获取下注记录
   * 每次传开始时间到服务器查询，
   * 从返回结果中存储最大的数据产生时间为下一次的起始时间
   */
  static async GetLog(client, req : GetLogReq) : Promise<[DocBetLogList,any]> {
    return await client.send("ForCustomerRpc/GetLog", req)
  }
}

/** 接入者提供 */
export class CustomerOfferedRpc {
  /** 修改单个玩家金币 */
  static async ModifyGold(client, req : ModifyGoldUidReq) : Promise<[GetBalanceUidResp,any]> {
    return await client.send("CustomerOfferedRpc/ModifyGold", req)
  }
  /** 获取单个玩家余额 */
  static async GetBalance(client, req : GetBalanceUidReq) : Promise<[GetBalanceUidResp,any]> {
    return await client.send("CustomerOfferedRpc/GetBalance", req)
  }
}