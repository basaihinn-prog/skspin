/* eslint-disable */
import type { Empty } from "../empty";


export interface SlotsPlayerRotateData {
}

export interface SlotsPlayerRotateData_Doc {
  /** @gotags: bson:"_id" */
  ID?: string;
  RotateCount?: number;
}

export interface SlotsPlayerRotateData_IncRotateReq {
  Pid?: number;
  Game?: string;
}

export interface SlotsPlayerRotateData_GetRotateReq {
  Pid?: number;
  Game?: string;
}

export interface SlotsPlayerRotateData_GetRotateResp {
  Count?: number;
}

export class SlotsPlayerRotateDataRpc {
  static async IncRotate(client, req : SlotsPlayerRotateData_IncRotateReq) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsPlayerRotateDataRpc/IncRotate", req)
  }
  static async GetRotate(client, req : SlotsPlayerRotateData_GetRotateReq) : Promise<[SlotsPlayerRotateData_GetRotateResp,any]> {
    return await client.send("slots.SlotsPlayerRotateDataRpc/GetRotate", req)
  }
}