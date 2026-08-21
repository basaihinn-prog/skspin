/* eslint-disable */
import type { Empty } from "../empty";


export enum PoolType {
  /** Invalid - 无效 */
  Invalid = 0,
  /** Slots - 所有slots */
  Slots = 1,
  /** BaiRen - 所有百人 */
  BaiRen = 2,
}

export interface IncSelfSlotsPoolReq {
  Pid?: number | undefined;
  Gold?: number | undefined;
  Type?: PoolType | undefined;
}

export interface GetSelfSlotsPoolReq {
  Pid?: number | undefined;
  Type?: PoolType | undefined;
}

export interface GetSelfSlotsPoolResp {
  Gold?: number | undefined;
}

export interface SetSelfSlotsPoolReq {
  Pid?: number | undefined;
  Gold?: number | undefined;
  Type?: PoolType | undefined;
}

export class SlotsPoolRpc {
  /** 更新个人奖池 */
  static async IncSelfSlotsPool(client, req : IncSelfSlotsPoolReq) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsPoolRpc/IncSelfSlotsPool", req)
  }
  /** 获取个人奖池 */
  static async GetSelfSlotsPool(client, req : GetSelfSlotsPoolReq) : Promise<[GetSelfSlotsPoolResp,any]> {
    return await client.send("slots.SlotsPoolRpc/GetSelfSlotsPool", req)
  }
  /** 设置个人奖池 */
  static async SetSelfSlotsPool(client, req : SetSelfSlotsPoolReq) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsPoolRpc/SetSelfSlotsPool", req)
  }
}