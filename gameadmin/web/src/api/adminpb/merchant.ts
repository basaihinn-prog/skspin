import {Empty} from "@/api/empty";
import {CommPageRequest} from "@/api/comm";
import {AddMaintenance} from "@/api/adminpb/info";

export interface MerchantGameReq {
    Gametype    : number,
    GameName    : string,
    GameOn      : number,
    UserName    : string
    PageIndex   : number
    PageSize    : number
    GameManufacturer: Array<string>
}

export interface MerchantGameListRep {
    CountAll: number,
    List: MerchantGameInterface[]
}

export interface MerchantGameInterface {
    GameId: string,               // 游戏编号
    Gametype: number,             // 游戏分类
    GameName: string,             // 游戏名称
    GameManufacturer: string,             // 游戏名称
    ConfigPath: string,           // 配置文件
    RTP: number,                  // RTP设置
    StopLoss: number,             // 止盈止损开关
    MaxMultipleOff: number,
    MaxMultiple   : number,       // 赢取最高押注倍数
    BetBase: string,              // 游戏投注
    GamePattern: number,          // 游戏类型
    Preset: number,               // 预设面额
    GameOn: number,               // 游戏状态
    FreeGameOff: number,          // 免费游戏开关
    ShowNameAndTimeOff: number,   // 商户是否显示游戏名称
    ShowNameAndTime: number,   // 游戏是否显示游戏名称
    ShowBagOff: number,             // 商户是是否显示背包
    ShowBag: number,             // 游戏是否显示背包
    GameDemo: string,             // 试玩链接
    DefaultCs: number,            // 默认投注
    DefaultLevel: number,         // 默认等级
    OnlineUpNum: number,
    OnlineDownNum: number,
    MaxWinPoints: number                     // 游戏投注
    BuyRTP: number                     // 游戏投注
}


export interface updateGameConfigParams {
    AppID?: string                    // 商户名称
    UserName?: string                    // 商户名称
    GameOn: number                      // 商户名称
    GameId: string                      // 游戏编号
    Preset: number                      // 预设面额
    StopLoss: number                    // 止盈止损开关
    GamePattern: number                 // 游戏类型
    BuyRTP:number
    ShowBag: number                 // 游戏类型
    ShowNameAndTimeOff: number                 // 游戏类型
    FreeGameOff: number                 // 免费游戏开关
    RTP: number                         // RTP设置
    MaxMultiple: number                 // 赢取最高押注倍数
    BetBase: string                     // 游戏投注
    MaxWinPoints: number                     // 游戏投注
    Gametype: string,
    GameName: string,
    GameManufacturer: string,
    OnlineUpNum: number,
    OnlineDownNum: number,
    ConfigPath: string,
    DefaultCs: number
    DefaultLevel: number
    MaxMultipleOff:number

}



export interface BatchEditGameRTPParam {
    AppID?: string                    // 商户名称
    GameList: string[]
    GamePattern: number
    RTP:      number
    MaxWinPoints:      number
    MaxMultiple:      number
    BuyRTP:number,
    BetMultiples: number
}


export interface BatchEditGameBetParam {
    AppID        ?:     string                    // 商户名称
    Manufacturer :      string[]
    GameIds      :      string
    OnlineUpNum  :      number
    OnlineDownNum:      number
    BetBase      :      string
}


export interface ReviewListParams extends CommPageRequest{
    AppID       :string
    ParentAppID :string
    StartTime   :number
    EndTime     :number

}
export interface ReviewListRes{
    All: number
    List: Review[]
}
export interface Review extends AddMaintenance{
    // 线路商审核状态
    ReviewStatus: number
    // 审核人
    Reviewer: string
}

export interface ReviewStatusParams {
    OperatorID: number
}



export interface RTPProtectListParams extends CommPageRequest{

    AppID        :string  // AppId
    IsProtection :number  // 开关`
    ProtectionRewardPercentLess: number
    ProtectionRotateCount: number

}

export interface RTPProtectListRep {
    CountAll: number,
    List: RTPProtectInfo[]
}

export interface RTPProtectInfo {
    AppID :string
    // 保护开关
    IsProtection :number
    // 保护次数
    ProtectionRotateCount :number
    // 保护内获取比率
    ProtectionRewardPercentLess :number

}


export class Merchant{

    // 获取商户游戏列表
    static async GetMerchantGames(client, req: MerchantGameReq): Promise<[MerchantGameListRep, any]> {
        return await client.send("AdminInfo/GameConfig", req)
    }
    // 更新游戏配置
    static async UpdateGameConfig(client, req: updateGameConfigParams): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateGameConfig", req)
    }
    // 更新游戏配置
    static async BatchSetGameOn(client, req: updateGameConfigParams): Promise<[Empty, any]> {
        return await client.send("AdminInfo/BatchSetGameOn", req)
    }
    // 批量修改RTP
    static async BatchEditGameRTP(client, req: BatchEditGameRTPParam): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateBatchGameRTP", req)
    }

    // 批量修改BET
    static async BatchSetGameBet(client, req: BatchEditGameBetParam): Promise<[Empty, any]> {
        return await client.send("AdminInfo/BatchSetGameBet", req)
    }
    static async BatchSetGameBetMultiples(client, req: BatchEditGameBetParam): Promise<[Empty, any]> {
        return await client.send("AdminInfo/BatchSetGameBetMult", req)
    }
    // 审批商户查询列表
    static async ApprovalOperatorList(client, req: ReviewListParams): Promise<[ReviewListRes, any]> {
        return await client.send("AdminInfo/GetMerchantReviewList", req)
    }

    // 审批商户
    static async ApprovalOperator(client, req: ReviewStatusParams): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateMerchantReviewStatus", req)
    }


    // 商户RTP保护
    static async RTPProtectList(client, req: RTPProtectListParams): Promise<[RTPProtectListRep, any]> {
        return await client.send("AdminInfo/GetExtraRtpConfig", req)
    }
    // 商户RTP保护
    static async SetRTPProtect(client, req: RTPProtectListParams): Promise<[Empty, any]> {
        return await client.send("AdminInfo/SetExtraRtpConfig", req)
    }



}
