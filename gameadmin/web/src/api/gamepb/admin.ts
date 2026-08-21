/* eslint-disable */
import type {Empty} from "../empty";
import type {DocBetLog, Game, GameListResp} from "./customer";
import {GameStatus, GameType} from "./customer";
import {CommPageRequest} from "@/api/comm";


export enum ModifyStatus {
    OK = 0,
    Err = 1,
}

export interface App {
    /** @gotags:bson:"_id" */
    AppID?: string;
    AppSecret?: string;
    /** 请求地址 */
    Address?: string;
    /** 备注 */
    Comment?: string;
    /**  */
    CreateAt?: string;
    Pid :string;
    Uid? : string;
    CurrencyKey? : string
    CurrencyName? : string
    LoginAt? : string
    TypeInfo? : string
    Bet? : string
    Win? : string
}

export interface AppListResp {
    List?: App[];
    AllCount:number
}
export interface CommResp {
    NoExitPId: string
}

export interface PlayerInfoReq {
    Pid?: number;
    /** AppID-Uid 需要传入AppID */
    Uid?: string;
}

export interface GetPlayerBalanceReq {
    Pid: number;
}

export interface GetBetDetailsUrlReq{
    Gid:string;
    BetID:string;
    Lang:string;
    Token:string;
}

export interface DocPlayer {
    /** 内部ID */
    Pid?: number;
    /** 外部id */
    Uid?: string;
    /** 所属产品 */
    AppID?: string;
    /** 最后登录时间 */
    LoginAt?: string;
    /** 创建时间 */
    CreateAt?: string;
    Bet?: number;
    Win?: number;
}

export interface DocModifyGoldLog {
    /** @gotags:bson:"_id" */
    ID?: string;
    /** 内部id */
    Pid?: number;
    /** 金币变化 */
    Change?: number;
    /** 修改后余额 */
    Balance?: number;
    /** 备注 */
    Comment?: string;
    /** 状态 */
    Status?: ModifyStatus;
    /** 请求时间 */
    ReqTime?: string;
    /** 回应时间 */
    RespTime?: string;
}

export interface PlayerListReq {
    AppID?: string;
    PageIndex?: number;
    PageSize?: number;
}

export interface PlayerListResp {
    Count?: number;
    List?: DocPlayer[];
}

export interface BHdownloadResp {
    ID: string;
}

export interface BetLogListReq {
    OperatorId?: number,
    OrderId?: string;
    Pid?: number;
    StartTime?: number;
    EndTime?: number;
    times?: []
    GameID?: string;
    Bet?: number;
    Win?: number;
    WinLose?: number;
    Balance?: number;
    LastId?: string;
}

export interface BetLogListResp {
    Count?: number;
    List?: DocBetLog[];
}

export interface ModifyLogListReq {
    /** 玩家ID，空=所有玩家 */
    Pid?: number;
    PageIndex?: number;
    PageSize?: number;
    TimeRange?: string[];
}
export interface GetOpenDaysByYearReq {
    Year?: string;
}

export interface ModifyLogListResp {
    Count?: number;
    List?: DocModifyGoldLog[];
}

export interface NewAppListReq {
    /** 产品列表 */
    // PageIndex: number;
    // PageSize: number;
    Pid?: string
    Uid?: string
}

export interface RestrictionParams {
    /** 产品列表 */
    Pids?: any[] | string
    Pid?: number
    RestrictionsStatus?: number
    RestrictionsMaxWin?: number
    RestrictionsMaxMulti?: number
}
export interface GameReq extends Game{
    ID?: string;
    Name?: number;
    BuyBetMulti?: string;
    Type?: GameType;
    Status?: GameStatus;
    Icon?: string;

}
export interface MaintenanceLogResponse {
    List?:MaintenanceLog[];
    Count:number
}

export interface MaintenanceLog {
    ID                : string
    MaintenanceMod    ?: number;
    MaintenanceStatus ?: number;
    StartTime         ?: number;
    EndTime           ?: number;
    RealEndTime       ?: number;
    OperatorId        ?: string;
    OperatorTime      ?: number;
}


export interface SystemConfig{
    HalfMaintenance: number
    EntireMaintenance: number
    Maintenance?: number
    StartTime?: number
    EndTime?: number
}

export interface GetExcleBetLogAddReq {

}
/** @ts: prefix(gamecenter) */
export class AdminGameCenter {
    /** 添加一个产品 */
    static async AddApp(client, req: App): Promise<[App, any]> {
        return await client.send("gamecenter/AdminGameCenter/AddApp", req)
    }

    /*static async AppList(client, req : Empty) : Promise<[AppListResp,any]> {
      return await client.send("gamecenter/AdminGameCenter/AppList", req)
    }*/
    /** 产品列表 */
    static async NewAppList(client, req: NewAppListReq): Promise<[AppListResp, any]> {
        return await client.send("AdminInfo/GetPlayerList", req)
    }
   /** 获取终身体限值 */
    static async GetPlayerRestrictionsList(client, req: NewAppListReq): Promise<[AppListResp, any]> {
        return await client.send("AdminInfo/GetPlayerRestrictionsList", req)
    }

   /** 获取终身体限值 */
    static async SetPlayerRestrictions(client, req: RestrictionParams): Promise<[CommResp, any]> {
        return await client.send("AdminInfo/SetPlayerRestrictions", req)
    }
   /** 获取终身体限值 */
    static async CancelPlayerRestrictions(client, req: RestrictionParams): Promise<[AppListResp, any]> {
        return await client.send("AdminInfo/CancelPlayerRestrictions", req)
    }

    /** 修改产品信息，Address/Comment 可修改 */
    static async ModifyApp(client, req: App): Promise<[Empty, any]> {
        return await client.send("gamecenter/AdminGameCenter/ModifyApp", req)
    }

    /** 获取用户详情 */
    static async PlayerInfo(client, req: PlayerInfoReq): Promise<[DocPlayer, any]> {
        return await client.send("gamecenter/AdminGameCenter/PlayerInfo", req)
    }

    /**  */
    static async PlayerList(client, req: PlayerListReq): Promise<[PlayerListResp, any]> {
        return await client.send("gamecenter/AdminGameCenter/PlayerList", req)
    }

    /** 下注日志 */
    static async BetLogList(client, req: BetLogListReq): Promise<[BetLogListResp, any]> {
        return await client.send("AdminInfo/GetBetLogListNew", req)
    }

    /** 金币修改日志 */
    static async ModifyLogList(client, req: ModifyLogListReq): Promise<[ModifyLogListResp, any]> {
        return await client.send("gamecenter/AdminGameCenter/ModifyLogList", req)
    }

    /** 游戏列表 */
    static async GameList(client, req: Empty): Promise<[GameListResp, any]> {
        return await client.send("AdminInfo/GetGameList", req)
    }
    /** 游戏列表--Operator */
    static async GameListOperator(client, req: Empty): Promise<[GameListResp, any]> {
        return await client.send("AdminInfo/GetGameListExc", req)
    }
    /** 修改一个游戏 */
    static async ModifyGame(client, req: Game): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateGame", req)
    }

    /** 添加一个游戏 */
    static async AddGame(client, req: Game): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddGame", req)
    }


    static async GameBetWinDataList(client, req: Game): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetGameBetWinData", req)
    }

    static async GetOpenDaysByYear(client, req: GetOpenDaysByYearReq): Promise<[Empty, any]> {
        return await client.send("mq/lotter/admin/getOpenDaysByYear", req)
    }

    static async GetPlayerBalance(client, req: GetPlayerBalanceReq): Promise<[GetPlayerBalanceResp, any]> {
        return await client.send("mq/gamecenter/player/balance", req)
    }
    static async GetPlrLotteryInfo(client, req: GetPlayerBalanceReq): Promise<[GetPlayerBalanceResp, any]> {
        // return await client.send("mq/lotter/admin/plrLotteryInfo", req)
        return new Promise(null)
    }
    static async GetBetDetailsUrl(client, req: GetBetDetailsUrlReq): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("mq/pggateway/getBetDetailsUrl", req)
    }

    static async Clear(client, req: GetBetDetailsUrlReq): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("/AdminInfo/clearUri", req)
    }

    static async GetPPBetDetailsUrl(client, req: GetBetDetailsUrlReq): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("mq/ppgateway/getBetDetailsUrl", req)
    }
    static async GetJiliBetDetailsUrl(client, req: GetBetDetailsUrlReq): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("mq/jiligateway/getBetDetailsUrl", req)
    }

    static async GetZIPBetLogAdd(client, req: BetLogListReq): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("AdminInfo/ZIPBetLogAdd", req)
    }

    static async GetExcleBetLogList(client, req: BetLogListReq): Promise<[GetExcleBetLogListResp, any]> {
        return await client.send("AdminInfo/ZIPBetLogList", req)
    }

    static async GetBHdownload(client, req: BHdownloadResp): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("AdminInfo/BHdownload", req)
    }

    // 查询对应商户的配置
    static async SelectMerchant(client, req: Empty): Promise<[SystemConfig, any]> {
        return await client.send("AdminInfo/getConfig", req)
    }
    // 修改添加
    static async ModifyConfig(client, req: SystemConfig): Promise<[GetBetDetailsUrlResp, any]> {
        return await client.send("AdminInfo/editConfig", req)
    }

    // 获取配置信息
    static async GetMaintenanceLog(client, req: CommPageRequest): Promise<[MaintenanceLogResponse, any]> {
        return await client.send("AdminInfo/getMaintenanceLog", req)
    }

    // 获取配置信息
    static async UploadGameFile(client, req: FormData): Promise<[MaintenanceLogResponse, any]> {
        return await client.send("AdminInfo/uploadGame", req)
    }


}
export interface GetBetDetailsUrlResp {
    Url?: string;
}
export interface GetExcleBetLogListResp {
    AllCount?: number;
    list?:any[]
}
export interface GetPlayerBalanceResp {
    Balance?:any
    Unsettled?:any
}
