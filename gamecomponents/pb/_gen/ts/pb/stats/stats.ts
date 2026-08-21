/* eslint-disable */
import type { SlotsStatsBuyGame_Doc, ThreeMonkeyRunTimeResp } from "../slots/data";


export interface JsonData {
  Data?: string | undefined;
}

export interface GetLast10DayRunTimeDataReq {
  Game?: string | undefined;
}

export interface GetLast10DayRunTimeDataResp {
  General?: ThreeMonkeyRunTimeResp | undefined;
  BuyGame?: SlotsStatsBuyGame_Doc[] | undefined;
}

export interface GetTotalExpenseReq {
  Days?: number | undefined;
  Game?: string | undefined;
}

export interface GetTotalExpenseResp {
  /** 总流水 */
  TotalExpense?:
    | number
    | undefined;
  /** 总产出 */
  TotalIncome?: number | undefined;
}

/** @ts: prefix(stats) */
export class AdminStatsRpc {
  static async GetLast10DayRunTimeData(client, req : GetLast10DayRunTimeDataReq) : Promise<[JsonData,any]> {
    return await client.send("stats/stats.AdminStatsRpc/GetLast10DayRunTimeData", req)
  }
  /** 获取最近n天的总流水与总产出 */
  static async GetTotalExpense(client, req : GetTotalExpenseReq) : Promise<[GetTotalExpenseResp,any]> {
    return await client.send("stats/stats.AdminStatsRpc/GetTotalExpense", req)
  }
}