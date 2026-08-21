/* eslint-disable */
import type { Empty } from "../empty";
import type { SlotsStatsBuyGame_Doc, ThreeMonkeyRunTimeResp } from "../slots/data";


export interface JsonData {
  Data?: string | undefined;
}

export interface AdminRecoverPanReq {
  UniqueStr?: string | undefined;
}

export interface GetLast10DayRunTimeDataResp {
  General?: ThreeMonkeyRunTimeResp | undefined;
  BuyGame?: SlotsStatsBuyGame_Doc[] | undefined;
}

export class AdminSlotsRpc {
  static async AdminRecoverPan(client, req : AdminRecoverPanReq) : Promise<[JsonData,any]> {
    return await client.send("adminslots.AdminSlotsRpc/AdminRecoverPan", req)
  }
  static async GenCombineData(client, req : Empty) : Promise<[Empty,any]> {
    return await client.send("adminslots.AdminSlotsRpc/GenCombineData", req)
  }
  static async GenerateData(client, req : Empty) : Promise<[Empty,any]> {
    return await client.send("adminslots.AdminSlotsRpc/GenerateData", req)
  }
  static async GetLast10DayRunTimeData(client, req : Empty) : Promise<[JsonData,any]> {
    return await client.send("adminslots.AdminSlotsRpc/GetLast10DayRunTimeData", req)
  }
  static async GetPoolStatus(client, req : Empty) : Promise<[JsonData,any]> {
    return await client.send("adminslots.AdminSlotsRpc/GetPoolStatus", req)
  }
  static async ExtractAction(client, req : JsonData) : Promise<[Empty,any]> {
    return await client.send("adminslots.AdminSlotsRpc/ExtractAction", req)
  }
  static async SwitchAutoCreateData(client, req : JsonData) : Promise<[Empty,any]> {
    return await client.send("adminslots.AdminSlotsRpc/SwitchAutoCreateData", req)
  }
}