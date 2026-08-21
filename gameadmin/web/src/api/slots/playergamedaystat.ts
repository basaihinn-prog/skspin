/* eslint-disable */
import type { Empty } from "../empty";


export interface SlotsPlayerGameDayStatRecord {
  Game?: string;
  Pid?: number;
  Date?: string;
  /** 特殊的Flag就是针对All，代表 */
  Flag?: string;
  /** 押注 */
  Flow?: number;
  /** 赢分 */
  Win?: number;
  /** 业绩 */
  YeJi?: number;
  /** 净输赢 */
  NetProfit?: number;
  /** 抽水 */
  Fee?: number;
  /** 个人奖池爆奖励 */
  SelfPoolReward?: number;
}

export interface SlotsDayBetReq {
  Channel?: string;
  Pid?: number;
  Game?: string;
  Start?: string;
  End?: string;
}

export interface SlotsGetAllDayInfoResp {
  List?: SlotsPlayerGameDayStatRecord[];
}

export interface GetGameDayWinCountReq {
  Game?: string;
  Day?: string;
}

export interface GetGameDayWinCountResp {
  Count?: number;
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