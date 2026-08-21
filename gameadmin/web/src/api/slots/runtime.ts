/* eslint-disable */
import type { Empty } from "../empty";
import type { ThreeMonkeyRunTimeData } from "./data";


export interface AddRunTimeDataReq {
  Pid?: number;
  Game?: string;
  Expense?: number;
  Income?: number;
  SelfPoolReward?: number;
  IsSmallGame?: boolean;
  IsFreeGameTurn?: boolean;
}

export interface GetTotalReq {
  Game?: string;
}

export interface GetTotalResp {
  Expense?: number;
  Income?: number;
}

export interface FixRunTimeDataReq {
  Day?: string;
  Game?: string;
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