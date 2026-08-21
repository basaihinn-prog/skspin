import {CommPageRequest} from "@/api/comm";


const GET_PLAYER_RETENTION_REPORT_DATA = "/AdminInfo/GetPlayerRetentionReportData"
const GET_PLAYER_RETENTION_GAME_REPORT_DATA = "/AdminInfo/GetPlayerRetentionGameReportData"
const GET_RANKING_LIST = "/AdminInfo/GetRankingList" // 获取排行榜

export interface GameUserRetainedReq extends CommPageRequest{
    GameId?:string
    Manufacturer?:string
    Operator?:number
    StartTime?:number
    EndTime?:number
    ReportId?:string
}

export interface GetRankingListParam extends CommPageRequest{
    OperatorId      ?:number //商户ID
    StartTime       ?:number
    EndTime         ?:number
    GameID          ?:string //游戏ID
    Manufacturer    ?:string
    RankingListType ?:number
    GameType ?:number
}

export interface RetainedListRes{
    All:number
    List:RetainedList[]
}
export interface RetainedList{
    ID                          : string
    RetentionPlayerCount        : number
    RetentionPlayer1d           : number
    RetentionPlayer3d           : number
    RetentionPlayer7d           : number
    RetentionPlayer14d          : number
    RetentionPlayer30d          : number
    Date                        : string
    AppID                       : string
}

export interface RankingListRes{
    Count:number
    List:RankingList[]
}
export interface RankingList{
    Pid          : number
    TotalWinLoss : number
    TotalWin     : number
    Balance      : number
    SpinCount    : number
    TotalBet     : number
    AppID        : string
    WinRate      : number
    CurrencyCode : string
    CurrencyName : string
    LoginAt      : string
    CreateAt     : string
}

export interface RetainedInfoRes{
    All:number
    List:RetainedGameList[]
}


export interface RetainedGameList{
    ID                     : string
    RetentionPlayerCount   : number
    RetentionPlayer1d      : number
    RetentionPlayer3d      : number
    RetentionPlayer7d      : number
    RetentionPlayer14d     : number
    RetentionPlayer30d     : number
    Date                   : string
    GameID                 : string
    GameName               : string
    GameGameManufacturer   : string
    GameIcon               : string

}





export class GameUserRetained{

    // 获取列表
    static async GameUserRetainedList(client, req: GameUserRetainedReq): Promise<[RetainedListRes, any]> {
        return await client.send(GET_PLAYER_RETENTION_REPORT_DATA, req)
    }
    // 获取详情
    static async GameUserRetainedInfo(client, req: GameUserRetainedReq): Promise<[RetainedInfoRes, any]> {
        return await client.send(GET_PLAYER_RETENTION_GAME_REPORT_DATA, req)
    }
    // 获取详情
    static async GetRankingList(client, req: GetRankingListParam): Promise<[RankingListRes, any]> {
        return await client.send(GET_RANKING_LIST, req)
    }
}
