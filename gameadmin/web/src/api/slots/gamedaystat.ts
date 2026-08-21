/* eslint-disable */
import type { Empty } from "../empty";


export interface SlotsGameDayStatRecord {
  Game?: string;
  Date?: string;
  /** 特殊的Flag就是针对All，代表 */
  Flag?: string;
  EnterCount?: number;
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

export interface SlotsGameDayStatDocsReq {
  StartDate?: string;
  EndDate?: string;
  Game?: string;
  Flags?: string[];
}

export interface SlotsGameDayStatDocsResp {
  List?: SlotsGameDayStatRecord[];
}

export class SlotsGameDayStatRpc {
  static async IncBoth(client, req : SlotsGameDayStatRecord) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsGameDayStatRpc/IncBoth", req)
  }
  static async IncAll(client, req : SlotsGameDayStatRecord) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsGameDayStatRpc/IncAll", req)
  }
  static async Docs(client, req : SlotsGameDayStatDocsReq) : Promise<[SlotsGameDayStatDocsResp,any]> {
    return await client.send("slots.SlotsGameDayStatRpc/Docs", req)
  }
}