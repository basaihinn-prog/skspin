/* eslint-disable */
import type { Empty } from "../empty";
import type { ThreeMonkeyRunTimeData } from "./data";


export interface AddRunTimeDataReq {
  Pid?: number | undefined;
  Game?: string | undefined;
  Expense?: number | undefined;
  Income?: number | undefined;
  SelfPoolReward?: number | undefined;
  IsSmallGame?: boolean | undefined;
  IsFreeGameTurn?: boolean | undefined;
}

export interface GetTotalReq {
  Game?: string | undefined;
}

export interface GetTotalResp {
  Expense?: number | undefined;
  Income?: number | undefined;
}

export interface FixRunTimeDataReq {
  Day?: string | undefined;
  Game?: string | undefined;
}

export class SlotsRuntimeRpc {
  /** 更新运行时数据 */
  static async AddRunTimeData(client, req : AddRunTimeDataReq) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsRuntimeRpc/AddRunTimeData", req)
  }
  static async GetTotal(client, req : GetTotalReq) : Promise<[GetTotalResp,any]> {
    return await client.send("slots.SlotsRuntimeRpc/GetTotal", req)
  }
  static async FixRunTimeData(client, req : FixRunTimeDataReq) : Promise<[ThreeMonkeyRunTimeData,any]> {
    return await client.send("slots.SlotsRuntimeRpc/FixRunTimeData", req)
  }
}