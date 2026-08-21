/* eslint-disable */


export interface GetYeJiRatioReq {
  Game?: GetYeJiRatioReq_GameType;
}

export enum GetYeJiRatioReq_GameType {
  Slots = 0,
}

export interface GetYeJiRatioResp {
  Ratio?: number;
}

export class SlotsConfigRpc {
  /** 获取业绩计算比例 */
  static async GetYeJiRatio(client, req : GetYeJiRatioReq) : Promise<[GetYeJiRatioResp,any]> {
    return await client.send("slots.SlotsConfigRpc/GetYeJiRatio", req)
  }
}