/* eslint-disable */
import type {SlotsStatsBuyGame_Doc, ThreeMonkeyRunTimeResp} from "../slots/data";
import {Game} from "@/api/gamepb/customer";
import {Empty} from "@/api/empty";


export interface JsonData {
    MonthData?: JsonDataListResp;
    SevenData?: JsonDataListResp;
    TotalData?: JsonDataListResp;
    WeekData?: JsonDataWeekResp;
}

export interface JsonDataListResp {
    BetAmount?: number;
    WinAmount?: number;
}

export interface JsonDataWeekResp {
    LastWeek?: JsonDataWeekListResp;
    NowWeek?: JsonDataWeekListResp;
}

export interface JsonDataWeekListResp {
    BetValue: number
    BuyBet: number
    BuyRateValue: number
    DauValue: number
    RateValue: number
    SystemValue: number
    Time: string
}

export interface GetGameEarningsData {
    Game?: string;
    Type?: number;
}

export interface GetawardPoolData {
}

export interface SetawardPoolData {
    AwardPool: Number
}

export interface GetLast10DayRunTimeDataResp {
    General?: ThreeMonkeyRunTimeResp;
    BuyGame?: SlotsStatsBuyGame_Doc[];
}

export interface GetTotalExpenseReq {
    Days?: number;
    Game?: string;
}

export interface GetTotalExpenseResp {
    /** 总流水 */
    TotalExpense?: number;
    /** 总产出 */
    TotalIncome?: number;
}

export interface GetBetWinTotalReq {
    /** 总流水 */
    Date?: string;
    /** 总产出 */
    OperatorId?: number;
}

export interface GetGamesLoginListReq {
    OperatorId: number,
    GameID: string,
    StartTime: number,
    EndTime: number,
    PageSize: number,
    PageNumber: number,
    Count: number
}

export interface GetGamesLoginListResp {
    Count: number,
    PageNumber: number,
    PageSize: number,
    List: any[]
}

export interface GetBetWinTotalResp {
    LastMonth?: [];
    NowDay?: [];
    NowMonth?: [];
    SevenDay?: [];
}

export interface DownLoadResp {
    /** json字符串 */
    Json?: string;
}

export interface GetSlotPoolHistoryTypeResp {
    /** json字符串 */
    Types?: [];
}

export interface GetEnsureSlotPoolHistoryListReq {
    OperatorId: number,
    StartTime: number,
    EndTime: number,
    OpName: string,
    // Change: number,
    AnimUserPid: number,
    Type: number,
    EnsureStatus: number,
    EnsureStartTime: number,
    EnsureEndTime: number,
    PageSize: number,
    PageNumber: number,
    Count: number
}

export interface GetSlotWinLoseLimitListReq {
    OperatorId: number,
    StartTime: number,
    EndTime: number,
    Pid: string,
    EnsureName: number,
    EnsureStatus: number,
    EnsureStartTime: number,
    EnsureEndTime: number,
    PageSize: number,
    PageNumber: number,
    Count: number
}

export interface GetDoEnsureslotPoolHistoryReq {
    ID:string
    Remark:string
}

export interface GetDoEnsureslotPoolHistoryResp {

}

/** @ts: prefix(stats) */
export class AdminStatsRpc {
    static async GetGameEarningsMethod(client, req: GetGameEarningsData): Promise<[JsonData, any]> {
        return await client.send("AdminInfo/GetGameEarningsData", req)
    }

    static async GetawardPool(client, req: GetawardPoolData): Promise<[JsonData, any]> {
        return await client.send("/mq/hilo/awardPool/get", req)
    }

    static async SetawardPool(client, req: SetawardPoolData): Promise<[JsonData, any]> {
        return await client.send("/mq/hilo/awardPool/set", req)
    }

    /** 表格下载 */
    static async DownLoad(client, req: GetGameEarningsData): Promise<[DownLoadResp, any]> {
        return await client.send("AdminInfo/GetGameEarningsData/DownLoad", req)
    }

    /** 获取最近n天的总流水与总产出 */
    static async GetTotalExpense(client, req: GetTotalExpenseReq): Promise<[GetTotalExpenseResp, any]> {
        return await client.send("stats/stats.AdminStatsRpc/GetTotalExpense", req)
    }

    /*  获取游戏运行时（汇总）*/
    static async GetBetWinTotal(client, req: GetBetWinTotalReq): Promise<[GetBetWinTotalResp, any]> {
        return await client.send("AdminInfo/GetBetWinTotal", req)
    }

    static async GetGamesLoginList(client, req: GetGamesLoginListReq): Promise<[GetGamesLoginListResp, any]> {
        return await client.send("AdminInfo/GetGamesLoginList", req)
    }

    // 奖池修改记录审核的类型
    static async GetSlotPoolHistoryType(client, req: {}): Promise<[GetSlotPoolHistoryTypeResp, any]> {
        return await client.send("AdminInfo/slotPoolHistoryType", req)
    }

    // 待审批日志
    static async GetEnsureSlotPoolHistoryList(client, req: GetEnsureSlotPoolHistoryListReq): Promise<[GetGamesLoginListResp, any]> {
        return await client.send("AdminInfo/ensureSlotPoolHistoryList", req)
    }

    static async GetSlotWinLoseLimitList(client, req: GetSlotWinLoseLimitListReq): Promise<[GetGamesLoginListResp, any]> {
        return await client.send("AdminInfo/SlotWinLoseLimitList", req)
    }

    // 审核奖池修改的记录
    static async GetDoEnsureslotPoolHistory(client, req: GetDoEnsureslotPoolHistoryReq): Promise<[GetDoEnsureslotPoolHistoryResp, any]> {
        return await client.send("AdminInfo/doEnsureslotPoolHistory", req)
    }

    static async GetDoEnsureSlotWinLoseLimit(client, req: GetDoEnsureslotPoolHistoryReq): Promise<[GetDoEnsureslotPoolHistoryResp, any]> {
        return await client.send("AdminInfo/doEnsureSlotWinLoseLimit", req)
    }


}
