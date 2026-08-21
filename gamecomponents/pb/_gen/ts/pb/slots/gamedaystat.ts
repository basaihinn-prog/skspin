/* eslint-disable */
import type { Empty } from "../empty";


export interface SlotsGameDayStatRecord {
  Game?: string | undefined;
  Date?:
    | string
    | undefined;
  /** 特殊的Flag就是针对All，代表 */
  Flag?: string | undefined;
  EnterCount?:
    | number
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

export interface SlotsGameDayStatDocsReq {
  StartDate?: string | undefined;
  EndDate?: string | undefined;
  Game?: string | undefined;
  Flags?: string[] | undefined;
}

export interface SlotsGameDayStatDocsResp {
  List?: SlotsGameDayStatRecord[] | undefined;
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