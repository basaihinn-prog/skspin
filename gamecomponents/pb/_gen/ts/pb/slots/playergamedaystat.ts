/* eslint-disable */
import type { Empty } from "../empty";


export interface SlotsPlayerGameDayStatRecord {
  Game?: string | undefined;
  Pid?: number | undefined;
  Date?:
    | string
    | undefined;
  /** 特殊的Flag就是针对All，代表 */
  Flag?:
    | string
    | undefined;
  /** 押注 */
  Flow?:
    | number
    | undefined;
  /** 赢分 */
  Win?:
    | number
    | undefined;
  /** 业绩 */
  YeJi?:
    | number
    | undefined;
  /** 净输赢 */
  NetProfit?:
    | number
    | undefined;
  /** 抽水 */
  Fee?:
    | number
    | undefined;
  /** 个人奖池爆奖励 */
  SelfPoolReward?: number | undefined;
}

export interface SlotsDayBetReq {
  Channel?: string | undefined;
  Pid?: number | undefined;
  Game?: string | undefined;
  Start?: string | undefined;
  End?: string | undefined;
}

export interface SlotsGetAllDayInfoResp {
  List?: SlotsPlayerGameDayStatRecord[] | undefined;
}

export interface GetGameDayWinCountReq {
  Game?: string | undefined;
  Day?: string | undefined;
}

export interface GetGameDayWinCountResp {
  Count?: number | undefined;
}

export class SlotsPlayerGameDayStatRpc {
  static async IncBoth(client, req : SlotsPlayerGameDayStatRecord) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsPlayerGameDayStatRpc/IncBoth", req)
  }
  static async IncAll(client, req : SlotsPlayerGameDayStatRecord) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsPlayerGameDayStatRpc/IncAll", req)
  }
  static async GetAllDayInfo(client, req : SlotsDayBetReq) : Promise<[SlotsGetAllDayInfoResp,any]> {
    return await client.send("slots.SlotsPlayerGameDayStatRpc/GetAllDayInfo", req)
  }
  static async GetPlrDayInfoByGame(client, req : SlotsDayBetReq) : Promise<[SlotsPlayerGameDayStatRecord,any]> {
    return await client.send("slots.SlotsPlayerGameDayStatRpc/GetPlrDayInfoByGame", req)
  }
  static async GetGameDayWinCount(client, req : GetGameDayWinCountReq) : Promise<[GetGameDayWinCountResp,any]> {
    return await client.send("slots.SlotsPlayerGameDayStatRpc/GetGameDayWinCount", req)
  }
}