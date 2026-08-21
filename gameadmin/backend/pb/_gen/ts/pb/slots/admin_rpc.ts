/* eslint-disable */
import type { Empty } from "../empty";
import type { GetSelfSlotsPoolReq, GetSelfSlotsPoolResp, SetSelfSlotsPoolReq } from "./pool";


/** @ts: prefix(slots) */
export class AdminSlotsRpc {
  /** 获取个人奖池 */
  static async GetSelfSlotsPool(client, req : GetSelfSlotsPoolReq) : Promise<[GetSelfSlotsPoolResp,any]> {
    return await client.send("slots/slots.AdminSlotsRpc/GetSelfSlotsPool", req)
  }
  /** 设置个人奖池 */
  static async SetSelfSlotsPool(client, req : SetSelfSlotsPoolReq) : Promise<[Empty,any]> {
    return await client.send("slots/slots.AdminSlotsRpc/SetSelfSlotsPool", req)
  }
}