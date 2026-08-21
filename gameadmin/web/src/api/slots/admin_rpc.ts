/* eslint-disable */
import type {Empty} from "../empty";
import type {GetSelfSlotsPoolReq, GetSelfSlotsPoolResp, SetSelfSlotsPoolReq} from "./pool";

export interface SetBalanceReq {
    Pid: number;
    AppID?: number;
    Balance: number;
    BeforeBalance: number;
}
export interface SetHistoryReq {
    Pid: number;
    Type: number;
}
export interface GetHistoryResp {
    AnimUserPid: number,
    Change: number,
    NewGold: number,
    OldGold: number,
    OpName: string,
    OpPid: number,
    Type: number,
    _id: string,
    time: string
}
/** @ts: prefix(slots) */
export class AdminSlotsRpc {
    /** 获取个人奖池 */
    static async GetSelfSlotsPool(client, req: GetSelfSlotsPoolReq): Promise<[GetSelfSlotsPoolResp, any]> {
        return await client.send("AdminInfo/GetSlotsPool", req)
    }

    /** 设置个人奖池 */
    static async SetSelfSlotsPool(client, req: SetSelfSlotsPoolReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/SetSlotsPool", req)
    }

    static async setBalance(client, req: SetBalanceReq): Promise<[Empty, any]> {
        return await client.send("mq/gamecenter/player/setBalance", req)
    }

    static async setHistory(client, req: SetHistoryReq): Promise<[GetHistoryResp, any]> {
        return await client.send("AdminInfo/getSlotPoolHistory", req)
    }

    static async setHistoryList(client, req: SetHistoryReq): Promise<[GetHistoryResp, any]> {
        return await client.send("mq/hilo/awardPool/historyList", req)
    }

}
