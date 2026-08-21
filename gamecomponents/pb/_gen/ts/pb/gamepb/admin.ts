/* eslint-disable */
import type { Empty } from "../empty";
import type { DocBetLog, Game, GameListResp } from "./customer";


export enum ModifyStatus {
  OK = 0,
  Err = 1,
}

export interface App {
  /** @gotags:bson:"_id" */
  AppID?: string | undefined;
  AppSecret?:
    | string
    | undefined;
  /** 请求地址 */
  Address?:
    | string
    | undefined;
  /** 备注 */
  Comment?:
    | string
    | undefined;
  /**  */
  CreateAt?: string | undefined;
}

export interface AppListResp {
  List?: App[] | undefined;
}

export interface PlayerInfoReq {
  Pid?:
    | number
    | undefined;
  /** AppID-Uid 需要传入AppID */
  Uid?: string | undefined;
}

export interface DocPlayer {
  /** 内部ID */
  Pid?:
    | number
    | undefined;
  /** 外部id */
  Uid?:
    | string
    | undefined;
  /** 所属产品 */
  AppID?:
    | string
    | undefined;
  /** 最后登录时间 */
  LoginAt?:
    | string
    | undefined;
  /** 创建时间 */
  CreateAt?: string | undefined;
  Bet?: number | undefined;
  Win?: number | undefined;
}

export interface DocModifyGoldLog {
  /** @gotags:bson:"_id" */
  ID?:
    | string
    | undefined;
  /** 内部id */
  Pid?:
    | number
    | undefined;
  /** 金币变化 */
  Change?:
    | number
    | undefined;
  /** 修改后余额 */
  Balance?:
    | number
    | undefined;
  /** 备注 */
  Comment?:
    | string
    | undefined;
  /** 状态 */
  Status?:
    | ModifyStatus
    | undefined;
  /** 请求时间 */
  ReqTime?:
    | string
    | undefined;
  /** 回应时间 */
  RespTime?: string | undefined;
}

export interface PlayerListReq {
  AppID?: string | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
}

export interface PlayerListResp {
  Count?: number | undefined;
  List?: DocPlayer[] | undefined;
}

export interface BetLogListReq {
  /** 查询平台，空=查询所有平台 */
  AppID?:
    | string
    | undefined;
  /** 查询平台，空=查询所有平台 */
  GameID?:
    | string
    | undefined;
  /** 玩家ID，空=所有玩家 */
  Pid?: number | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
  TimeRange?: string[] | undefined;
}

export interface BetLogListResp {
  Count?: number | undefined;
  List?: DocBetLog[] | undefined;
}

export interface ModifyLogListReq {
  /** 玩家ID，空=所有玩家 */
  Pid?: number | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
  TimeRange?: string[] | undefined;
}

export interface ModifyLogListResp {
  Count?: number | undefined;
  List?: DocModifyGoldLog[] | undefined;
}

/** @ts: prefix(gamecenter) */
export class AdminGameCenter {
  /** 添加一个产品 */
  static async AddApp(client, req : App) : Promise<[App,any]> {
    return await client.send("gamecenter/AdminGameCenter/AddApp", req)
  }
  /** 产品列表 */
  static async AppList(client, req : Empty) : Promise<[AppListResp,any]> {
    return await client.send("gamecenter/AdminGameCenter/AppList", req)
  }
  /** 修改产品信息，Address/Comment 可修改 */
  static async ModifyApp(client, req : App) : Promise<[Empty,any]> {
    return await client.send("gamecenter/AdminGameCenter/ModifyApp", req)
  }
  /** 获取用户详情 */
  static async PlayerInfo(client, req : PlayerInfoReq) : Promise<[DocPlayer,any]> {
    return await client.send("gamecenter/AdminGameCenter/PlayerInfo", req)
  }
  /**  */
  static async PlayerList(client, req : PlayerListReq) : Promise<[PlayerListResp,any]> {
    return await client.send("gamecenter/AdminGameCenter/PlayerList", req)
  }
  /** 下注日志 */
  static async BetLogList(client, req : BetLogListReq) : Promise<[BetLogListResp,any]> {
    return await client.send("gamecenter/AdminGameCenter/BetLogList", req)
  }
  /** 金币修改日志 */
  static async ModifyLogList(client, req : ModifyLogListReq) : Promise<[ModifyLogListResp,any]> {
    return await client.send("gamecenter/AdminGameCenter/ModifyLogList", req)
  }
  /** 游戏列表 */
  static async GameList(client, req : Empty) : Promise<[GameListResp,any]> {
    return await client.send("gamecenter/AdminGameCenter/GameList", req)
  }
  /** 修改一个游戏 */
  static async ModifyGame(client, req : Game) : Promise<[Empty,any]> {
    return await client.send("gamecenter/AdminGameCenter/ModifyGame", req)
  }
  /** 添加一个游戏 */
  static async AddGame(client, req : Game) : Promise<[Empty,any]> {
    return await client.send("gamecenter/AdminGameCenter/AddGame", req)
  }
}