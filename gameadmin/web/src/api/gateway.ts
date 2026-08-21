/* eslint-disable */


export interface OnlineReq {
  Pid?: number;
}

export interface OnlineResp {
}

export class General {
  static async Online(client, req : OnlineReq) : Promise<[OnlineResp,any]> {
    return await client.send("General/Online", req)
  }
  static async Offline(client, req : OnlineReq) : Promise<[OnlineResp,any]> {
    return await client.send("General/Offline", req)
  }
}