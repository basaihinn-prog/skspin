/* eslint-disable */
import type { Empty } from "../empty";
import exp from "constants";


export interface Log {
  /** 玩家ID */
  Pid?: number;
  /** 下注 */
  Bet?: number;
  /** 输赢（是指产出，不是净输赢） */
  Win?: number;
  /** 游戏ID，用于标识是哪个游戏 */
  GameID?: string;
  /** 对局ID，方便查询 */
  RoundID?: string;
  /** 注释 */
  Comment?: string;
}

export interface AddLogReq {
  List?: Log[];
}

export interface GetBalanceReq {
  Pid?: number;
}

export interface GetBalanceResp {
  Balance?: number;
}

export interface ModifyGoldReq {
  /** 玩家ID */
  Pid?: number;
  /** 金币变化：>0 加钱， <0 扣钱 */
  Change?: number;
  /** 注释 */
  Comment?: string;
}

export interface TokenReq {
  Token?: string;
}

export interface TokenResp {
  /** 返回websocket 服务地址 */
  WsUrl?: string;
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

export interface Manufacturer {
  Id               :string
  ManufacturerName :string
  ManufacturerCode :string
}

export interface ManufacturerListResponse{
  List :Manufacturer[]
}

// 获取厂商
export class Manufacturer{

  static async GetManufacturerList(client, req : Empty) : Promise<[ManufacturerListResponse,any]> {
    return await client.send("AdminInfo/manufacturerList", req)
  }
}
