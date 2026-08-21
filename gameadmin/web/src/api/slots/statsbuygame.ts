/* eslint-disable */
import type { Empty } from "../empty";
import type { SlotsStatsBuyGame_Doc, SlotsStatsBuyGame_QueryReq, SlotsStatsBuyGame_QueryResp } from "./data";


export class SlotsStatsBuyGameRpc {
  static async Add(client, req : SlotsStatsBuyGame_Doc) : Promise<[Empty,any]> {
    return await client.send("slots.SlotsStatsBuyGameRpc/Add", req)
  }
  static async Query(client, req : SlotsStatsBuyGame_QueryReq) : Promise<[SlotsStatsBuyGame_QueryResp,any]> {
    return await client.send("slots.SlotsStatsBuyGameRpc/Query", req)
  }
}