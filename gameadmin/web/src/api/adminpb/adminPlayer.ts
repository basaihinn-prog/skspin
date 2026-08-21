import {NewAppListReq} from "@/api/gamepb/admin";
import {GameTypeListRep} from "@/api/adminpb/game";
import OperatorMaintenance from "@/pages/systemManagement/merchants/operatorMaintenance.vue";
import {AddMaintenance} from "@/api/adminpb/info";
import {Empty} from "@/api/empty";
import {CommPageRequest} from "@/api/comm";

const ADMIN_PLAYER_INFO = "AdminInfo/GetPlayerInfo"
const CreatePlayerRTP = "AdminInfo/CreatePlayerRTP"
const CLOSE_PLAYER_RTP = "AdminInfo/CancelPlayerRTPControlByPId"
const PLAYER_RTP_CONTROL_LIST = "AdminInfo/PlayerRTPControlList"
const CANCEL_CONTROL_LIST = "AdminInfo/CancelPlayerRTPControl"

const PLAYER_PAY_LIST = "AdminInfo/GetPlayRTP"
const PLAYER_PAY_INSERT = "AdminInfo/InsertPlayRTP"


export interface PlayerReq {
    Uid?: string
    OperatorId?: number
    Pid?: number
}

export interface PlayerResponse {
    GameID?: string
    RTPControlStatus?: number
    Operator?: AddMaintenance
    OnlineOperator?: AddMaintenance
    BuyRTP:number
    OperatorId?: number
    Id?: string,
    Pid?: number,
    Uid?: string,
    AppID?: string,
    CurrencyKey?: string,
    CurrencyName?: string,
    LoginAt?: string,
    CreateAt?: string,
    TypeInfo?: BetInfo[],
    Bet?: number,
    Win?: number,
    Status?: string,
    PlatformPay?: number
    onlineAppID?: string
    onlinePlatformPay?: number
    totalWin?: number
    PlayerRTPControllist: any[]
}


interface BetInfo {
    SpinCnt?: number,
    Bet?: number,
    Win?: number,
    BuyGame?: number,
    BuyGameBet?: number,
    BuyGameWin?: number,
    AvgBet?: number
}

export interface PlayerRTP {

    GameID        ?:string  //游戏列表
    Pid           ?:number   //玩家ID
    AppID         ?:string  //所属商户
    AppIDandPlayerID         ?:string  //所属商户
    ContrllRTP    ?:number //控制RTP
    BuyRTP: number, //控制RTP
    AutoRemoveRTP ?:number //自动解除RTP
    PersonWinMaxScore?: number, //自动解除RTP
    PersonWinMaxMult?: number //自动解除RTP
}

export interface PlayerRTPListRequest {
    OperatorId?: number               // 商户ID
    Pid?: number                      // 用户登录账号
    ControlTimeStart?: string         // 控制时间开始
    ControlTimeEnd?: string           // 控制时间结束
    GameId?: string                   // 控制时间结束
    Page:number
    PageSize: number
}



export interface PlayerRTPListResponse{
    List: PlayerRTPData[]
    Count: number
}

export interface PlayerRTPData{
    ControlTime: string             //控制时间
    OperatorName: string            //商户名
    UId: string                     //用户账号
    GameName: string                //游戏名称
    GameRTP: number                 //游戏配置RTP
    PlayerHistoryRTP: number        //玩家历史RTP
    ControlRTP: number              //控制RTP
    ControllingRTP: number          //被控制时RTP
    AutoRemoveRTP: number           //自动解除RTP
}

export interface SelectPlayRTPRequest extends CommPageRequest{

    DateType   :number //是最后登录时间还是创建时间
    AppID      :string //商户
    StartTime  :string //起始时间
    EndTime    :string //结束时间
    NowBet     :number //今日下注
    historyBet :number //历史下注
    NowWin     :number //今日输赢
    HistoryWin :number //历史输赢

}




export interface PlayRTPListResponse{
    List: PlayRTP[]
    Count: number
}




export interface PlayRTP{
    CreateTime: string
    LastLoginTime: string
    AppID      : string
    UId        : string
    NowBet     : number //今日下注
    historyBet : number //历史下注
    NowWin     : number //今日输赢
    HistoryWin : number //历史输赢
    Balance    : number

}




export interface CancelControlRequest{
    Ids: string
}







export class AdminPlayer {
    static async GetPlayerInfo(client, req: PlayerReq): Promise<[PlayerResponse, any]> {
        return client.send(ADMIN_PLAYER_INFO, req)
    }
    static async EditPlayerRTP(client, req: PlayerRTP): Promise<[Empty, any]> {
        return client.send(CreatePlayerRTP, req)
    }
    static async ClosePlayerRTP(client, req: PlayerRTP): Promise<[Empty, any]> {
        return client.send(CLOSE_PLAYER_RTP, req)
    }
    static async GetPlayerRTPList(client, req: PlayerRTPListRequest): Promise<[PlayerRTPListResponse, any]> {
        return client.send(PLAYER_RTP_CONTROL_LIST, req)
    }
    static async CancelPlayerRTPControl(client, req: CancelControlRequest): Promise<[PlayerRTPListResponse, any]> {
        return client.send(CANCEL_CONTROL_LIST, req)
    }
    static async PlayerPayInfoList(client, req: SelectPlayRTPRequest): Promise<[PlayRTPListResponse, any]> {
        return client.send(PLAYER_PAY_LIST, req)
    }
    static async PlayerPayInfoInsert(client, req: CancelControlRequest): Promise<[PlayerRTPListResponse, any]> {
        return client.send(PLAYER_PAY_INSERT, req)
    }
}
