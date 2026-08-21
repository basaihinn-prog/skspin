/* eslint-disable */


import {AppListResp, NewAppListReq} from "@/api/gamepb/admin";

export enum SearchType {
    Day = 0,
    Mon = 1,
}

export interface PlayerAnalysis {
    /** @gotags: bson:"_id" */
    ID?: string;
    Pid?: number;
    /**  */
    Bet?: number;
    /**  */
    Win?: number;
    AppID: string;
    BetAmount: number
    Date: string;
    LoginPlrCount: number;
    RegistPlrCount: number;
    RoundCount: number;
    SpinCount: number;
    WinAmount: number;
}

export interface AppAnalysis {
    /** @gotags: bson:"_id" */
    ID?: string;
    AppID?: string;
    Bet?: number;
    Win?: number;
}

export interface GameAnalysis {
    /** @gotags: bson:"_id" */
    CurrencyName?: string,
    BetAmount?: number,
    WinAmount?: number,
    SpinCount?: number,
    SpinPlrCount?: number,
    EnterPlrCount?: number,
    Date?: string,
    Game?: string,
    AppID?: string,
    BetAmountROI?: number,
    WinAmountROI?: number,
    EnterPlrCountROI?: number,
    SpinPlrCountROI?: number,
    pointsROI?: number,
    profitabilityROI?: number,
}

export interface PlayerAnalysisReq {
    PageIndex: number,
    PageSize: number,
    Operator: number,
    Date: string,
    EndDate: string,
    Type: string,
}

export interface PokerOperatorReq {
    PageIndex: number,
    PageSize: number,
    Operator: number,
    Date: string,
    EndDate: string,
    Type: string,
}

export interface PlayerAnalysisResp {
    Count?: number;
    List?: PlayerAnalysis[];
}

export interface AppAnalysisReq {
    AppID?: string;
    Search?: SearchType;
    PageIndex?: number;
    PageSize?: number;
}

export interface AppAnalysisResp {
    Count?: number;
    List?: AppAnalysis[];
}

export interface GameAnalysisReq {
    Date?: string;
    Operator?: number;
    PageIndex?: number;
    PageSize?: number;
}

export interface GetGameDataListReq {
    Date?: string;
    EndDate?: string;
    Game?: string;
    Operator?: number;
    PageIndex?: number;
    PageSize?: number;
}
export interface GameAnalysisResp {
    Count?: number;
    List?: GameAnalysis[];
    YesterdayList?: GameAnalysis[];
}
export interface GetGameDataListResp {
    Count?: number;
    List?: GameAnalysis[];
}
export interface PlayerReportReq {
    Date: string,
    EndDate: string,
    Operator: number,
    PageIndex: number,
    PageSize: number
}

export interface PlayerReportResp {
    Count?: number;
    List?: PlayerReportRespList[];
}

export interface PlayerReportRespList {
    /** @gotags: bson:"_id" */
    Pid?: number;
    AppID: string;
    BetAmount: number
    Date: string;
    RoundCount: number;
    SpinCount: number;
    WinAmount: number;
}

/** @ts: prefix(gamecenter) */
export class AdminAnalysis {
    static async Player(client, req: PlayerAnalysisReq): Promise<[PlayerAnalysisResp, any]> {
        return await client.send("AdminInfo/GetOperatorReportData", req)
    }

    static async GetPokerOperator(client, req: PokerOperatorReq): Promise<[PlayerAnalysisResp, any]> {
        return await client.send("AdminInfo/GetPokerOperatorReportData", req)
    }

    /** 平台/玩家汇总 */
    static async PlayerReport(client, req: PlayerReportReq): Promise<[PlayerReportResp, any]> {
        return await client.send("AdminInfo/GetPlayerReportData", req)
    }

    static async App(client, req: AppAnalysisReq): Promise<[AppAnalysisResp, any]> {
        return await client.send("gamecenter/gc.AdminAnalysis/App", req)
    }

    static async Game(client, req: GameAnalysisReq): Promise<[GameAnalysisResp, any]> {
        return await client.send("AdminInfo/GetGameData", req)
    }
    static async GameAnalysisList(client, req: GetGameDataListReq): Promise<[GetGameDataListResp, any]> {
        return await client.send("AdminInfo/GetGameDataList", req)
    }
}
