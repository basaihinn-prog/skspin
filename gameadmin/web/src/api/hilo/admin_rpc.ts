/* eslint-disable */
import type { Empty } from "../empty";
import type { GetAwardPoolResp, RuntimeDataResp, SetAwardPoolReq } from "./msg";


/** @ts: prefix(Hilo) */
export class AdminHiloRpc {
  static async GetAwardPool(client, req : Empty) : Promise<[GetAwardPoolResp,any]> {
    return await client.send("Hilo/hilo.AdminHiloRpc/GetAwardPool", req)
  }
  static async SetAwardPool(client, req : SetAwardPoolReq) : Promise<[Empty,any]> {
    return await client.send("Hilo/hilo.AdminHiloRpc/SetAwardPool", req)
  }
  static async RuntimeData(client, req : Empty) : Promise<[RuntimeDataResp,any]> {
    return await client.send("Hilo/hilo.AdminHiloRpc/RuntimeData", req)
  }
}