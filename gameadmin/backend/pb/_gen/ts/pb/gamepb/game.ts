/* eslint-disable */
import type { Empty } from "../empty";


export interface Log {
  /** 玩家ID */
  Pid?:
    | number
    | undefined;
  /** 下注 */
  Bet?:
    | number
    | undefined;
  /** 输赢（是指产出，不是净输赢） */
  Win?:
    | number
    | undefined;
  /** 游戏ID，用于标识是哪个游戏 */
  GameID?:
    | string
    | undefined;
  /** 对局ID，方便查询 */
  RoundID?:
    | string
    | undefined;
  /** 注释 */
  Comment?:
    | string
    | undefined;
  /** 余额 */
  Balance?: number | undefined;
}

export interface AddLogReq {
  List?: Log[] | undefined;
}

export interface GetBalanceReq {
  Pid?: number | undefined;
}

export interface GetBalanceResp {
  Balance?:
    | number
    | undefined;
  /** 玩家外部ID */
  Uid?: string | undefined;
}

export interface ModifyGoldReq {
  /** 玩家ID */
  Pid?:
    | number
    | undefined;
  /** 金币变化：>0 加钱， <0 扣钱 */
  Change?:
    | number
    | undefined;
  /** 注释 */
  Comment?: string | undefined;
}

export interface TokenReq {
  Token?: string | undefined;
}

export interface TokenResp {
  /** 返回websocket 服务地址 */
  WsUrl?: string | undefined;
}

/** 给游戏调用 */
export class GameRpc {
  /** 金币操作 */
  static async ModifyGold(client, req : ModifyGoldReq) : Promise<[GetBalanceResp,any]> {
    return await client.send("GameRpc/ModifyGold", req)
  }
  /** 获取玩家余额 */
  static async GetBalance(client, req : GetBalanceReq) : Promise<[GetBalanceResp,any]> {
    return await client.send("GameRpc/GetBalance", req)
  }
  /** 写入下注日志 */
  static async AddLog(client, req : AddLogReq) : Promise<[Empty,any]> {
    return await client.send("GameRpc/AddLog", req)
  }
  /** 验证token */
  static async ValidateToken(client, req : TokenReq) : Promise<[Empty,any]> {
    return await client.send("GameRpc/ValidateToken", req)
  }
}