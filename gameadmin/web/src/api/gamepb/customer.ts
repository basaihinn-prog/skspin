/* eslint-disable */
import type {Empty} from "../empty";


export enum GameType {
    /** Slot - 电子 */
    Slot = 0,
    /** Fish - 捕鱼游戏 */
    MINI = 1,
    /** Poker - 百人 */
    Poker = 3,
    /** Poker - 彩票 */
    Lotter = 4,
}

export enum GameStatus {
    /** Open - 正常 */
    Open = 0,
    /** Maintenance - 维护中（列表中出现） */
    Maintenance = 1,
    /** Hide - 隐藏（列表中不出现） */
    Hide = 2,
    Close = 3,
}

export enum WebSideMap {
    FROM_SITE_SANOOK = 0,
    FROM_SITE_THAIRATH = 1,
    FROM_SITE_MYHORA = 2,
    FROM_SITE_GLO = 3,
}

export interface DocBetLog {
    /** @gotags:bson:"_id" */
    ID?: string;
    /** 内部玩家ID */
    Pid?: number;
    /** 外部玩家ID */
    UserID?: string;
    /** 游戏ID */
    GameID?: string;
    /** 对局ID，方便查询 */
    RoundID?: string;
    /** 下注 */
    Bet?: number;
    /** 输赢 */
    Win?: number;
    /** 注释 */
    Comment?: string;
    /** 数据插入时间 */
    InsertTime?: string;
    /**  */
    AppID?: string;
}

export interface DocBetLogList {
    List?: DocBetLog[];
}

export interface GetLogReq {
    From?: string;
}

export interface LoginReq {
    /** 玩家id */
    UserID?: string;
    /** 游戏ID */
    GameID?: string;
    /** 游戏语言 */
    Language?: string;
}

export interface LoginResp {
    /** 游戏启动链接 */
    Url?: string;
}

export interface Game {
    GameId?: string;
    Id?: string;
    LineNum?: number;
    ManufacturerName?: string
    BuyBetMulti?: string;
    Type?: GameType;
    ChangeBetOff: number,
    Status?: GameStatus;
    Icon?: string;
    GameNameConfig?:string
}

export interface GameListResp {
    List?: Game[];
}

export interface dataListResp {
    data?: any;
}

export interface ModifyGoldUidReq {
    /** 玩家ID */
    UserID?: string;
    /** 金币变化：>0 加钱， <0 扣钱 */
    Change?: number;
    /** 注释 */
    Comment?: string;
}

export interface GetBalanceUidReq {
    UserID?: string;
}

export interface GetBalanceUidResp {
    Balance?: number;
}

/** 给接入者调用 */
export class ForCustomerRpc {
    /** 获取游戏列表 */
    static async GameList(client, req: Empty): Promise<[GameListResp, any]> {
        return await client.send("ForCustomerRpc/GameList", req)
    }

    /** 获取游戏启链接 */
    static async Login(client, req: LoginReq): Promise<[LoginResp, any]> {
        return await client.send("ForCustomerRpc/Login", req)
    }

    /**
     * 获取下注记录
     * 每次传开始时间到服务器查询，
     * 从返回结果中存储最大的数据产生时间为下一次的起始时间
     */
    static async GetLog(client, req: GetLogReq): Promise<[DocBetLogList, any]> {
        return await client.send("ForCustomerRpc/GetLog", req)
    }
}

/** 接入者提供 */
export class CustomerOfferedRpc {
    /** 修改单个玩家金币 */
    static async ModifyGold(client, req: ModifyGoldUidReq): Promise<[GetBalanceUidResp, any]> {
        return await client.send("CustomerOfferedRpc/ModifyGold", req)
    }

    /** 获取单个玩家余额 */
    static async GetBalance(client, req: GetBalanceUidReq): Promise<[dataListResp, any]> {
        return await client.send("CustomerOfferedRpc/GetBalance", req)
    }
}
