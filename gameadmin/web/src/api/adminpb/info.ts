/* eslint-disable */
import type {Empty} from "../empty";
import type {AdminStatus, DBAdminer} from "./auth";


/** 权限组表 */
export interface DBAuth {
    /** @gotags:bson:"_id" */
    ID?: number;
    /** 菜单id */
    MenuIds?: string;
    /** 权限组名称 */
    Name?: string;
    /** 备注 */
    Remark?: string;
    /** 创建时间 */
    CreateTime?: string;
}

export interface AdminListResp {
    Count?: number;
    List?: DBAdminer[];
}

export interface AddAdminReq {
    Username?: string;
    GroupId?: number;
}

export interface AddOperatorAdminReq {
    UserName?: string;
    Password?: string;
    AppID?: string;
}

export interface AddMaintenance {
    UserName:string,                         //  商户主账号
    OperatorType: number,                    //  商户类型
    LineMerchantID?: string,                  //  商户类型
    Password?: number[],                  //  商户类型
    Remark: string,                          //  商户备注
    BelongingCountry?: string,                //  归属国家
    Name: string,                            //  上级ID
    AppID:string,                            //  商户名称
    PlatformPay:number,                      //  平台费
    TurnoverPay?: number,                      //  商户类型
    CooperationType:string,                  //  合作模式
    CurrencyCtrlStatus:number,                     //  合作模式
    Advance:number,                          //  预付款金额
    CurrencyKey:string,                      //  商户币种
    CurrencyKeyName:string,                  //  商户币种
    Contact: any,                            //  联系方式
    WalletMode:number,                       //  钱包类型
    Surname:string,                          //  会员前缀
    Lang:string,                             //  客户端默认语言
    ServiceIp:string,                        //  服务器地址
    WhiteIps:any,                            //  服务器IP白名单
    Address:string,                          //  服务器回调地址
    UserWhite:string,                        //  用户白名单
    LoginOff: number,                        //  服务器IP白名单
    FreeOff: number,                         //  服务器回调地址
    RTPOff: number,                          //  服务器回调地址
    DormancyOff: number,                     //  用户白名单
    RestoreOff: number,                      //  记录开启
    ManualFullScreenOff: number,             //  免游游戏开关板
    NewGameDefaulOff: number,                //  防休眠开启
    BuyRTPOff: number,                      //  购买游戏开关
    MassageOff: number,                      //  断复开关
    MassageIp: string,                       //  手动全屏开启
    PlayerRTPSettingOff: number              //
    StopLoss: number,                        //  消息推送
    MaxMultipleOff: number,                  //  消息推送地址
    Robot: string,                           //  商户类型
    ChatID: string,                          //  商户类型
    Balance: number,                         //  商户类型
    ReviewType: number,                      //  商户类型
    BalanceThreshold: number,                //  余额不足阈值
    BalanceThresholdInterval: number,        //  余额不足间隔、
    ArrearsThresholdInterval: number,        // 欠费提示间隔
    ReviewStatus: number,        // 欠费提示间隔
    _LineMerchant: number,                   //  游戏新型预设
}



export interface EditMaintenance extends AddMaintenance{
    NextRate:number                         // 下个月的费率
    Status:number,                          // 商户状态
    RTPOff:number,                          // 商户状态
    ShowNameAndTimeOff?:number,             // 商户状态
    ShowExitBtnOff?:number,                    // 商户状态
    ExitLink?:string,                       // 商户状态
    MaxWinPointsOff?: number,
    HighRTPOff?:number,                     // 超高RTP
    CurrencyVisibleOff?: number,                   // 是否显示币种
    DefaultManufacturerOn?:Array<string>,
    PGConfig?: PGConfig
    // JILI游戏配置
    JILIConfig?: JILIConfig
    // PP游戏配置
    PPConfig?: PPConfig
    TADAConfig?: JILIConfig
    AdminRobot?: string,                             //总控机器人的chatid
    AdminChatID?: string,                            //余额告警阈值
    BalanceAlert?: number,                            //余额告警时间 暂时固定为24每天0点进行检测并触发余额告警
    BalanceAlertTimeInterval?: number,
}



export interface AddOperatorUserReq {
    UserName?: string;
    Password?: string;
    AppID?: string;
}

export interface AddOperatorReq {
    name?: string;
    AppID?: string;
}

export interface AdminPid {
    Pid?: number;
}

export interface UpdateOpenGooleReq {
    Id?: number;
    Open?: boolean;
}

export interface AdminerStatusReq {
    Pid?: number;
    Status?: AdminStatus;
}

export interface GetOperatorListReq {
    AllCount?: number;
    Status?: DBMenu[];
    PageIndex?: number;
    PageSize?: number;
}

export interface GetUpdatePlayerStatusReq {
    Pid:number;
    Status:number;
    Uid?:string;
    AppID?:string;
}

export interface GetExcludeGameAllReq {
    ID?:string
    Name?:string
    Remark?:string
    ExcluedGameIds?:object
}
/** 菜单表 */
export interface DBMenu {
    /** @gotags:bson:"_id" */
    ID?: number;
    /** 菜单标题 */
    Title?: string;
    /** 菜单pid 0为顶级pid */
    Pid?: number;
    /** 节点 */
    Url?: string;
    /** 图标icon */
    Icon?: string;
    /** 排序 */
    Sort?: number;
}

export interface AdminInfoResp {
    ID?: number;
    Username?: string;
    PermissionId?: number;
    MenuList?: DBMenu[];
    OperatorId?: number;
    OperatorName:string;
    AppID:string;
    IsDefaultPermission:boolean
}

export interface UpdatePasswdReq {
    OldPasswd?: string;
    NewPasswd?: string;
    ConfirmPasswd?: string;
}

export interface DeleteGameReq {
    ID?: string;
}

export interface PageListReq {
    PageIndex?: number;
    PageSize?: number;
    Query?: string;
    Sort?: string;
}

export interface GetSaveMsg4GameUpdateMsgReq{
    Msg:string
}
export interface GetSaveOperator4GameUpdateMsgReq {
    DestOperatorID:[]
}
export interface GetConfig4GameUpdateMsgResp{
    Msg:string
    DestOperatorID:[]
}
export interface GetSendMsg4GameUpdateReq{
    Msg:string
    DestOperatorID:[]
}
export interface GetResult4GameUpdateMsgResp{
    ID:number,
    Completed:boolean,
    Ms:string,
    DestOperatorID:[],
    FailOperatorID: [],
    SuccessOperatorID: []
}
export interface GetStatus4GameUpdateMsgResp{
    Status:number
}


export interface UpdateGameUniversalConfigParam {
    AppID       : string
    PGConfig    : PGConfig
    JILIConfig  : JILIConfig
    PPConfig    : PPConfig
    TADAConfig  : JILIConfig
}


// pp游戏通用配置
export interface PPConfig {
    ShowNameAndTimeOff: number
    CurrencyVisibleOff: number
}

// pg游戏通用配置
export interface PGConfig {
    CarouselOff        : number
    ShowNameAndTimeOff : number
    CurrencyVisibleOff : number
    StopLoss           : number
    ExitBtnOff         : number
    ExitLink           : string
    OfficialVerify     : number
}

// jili游戏通用配置
export interface JILIConfig {
    BackPackOff        :number
    CurrencyVisibleOff :number
    OpenScreenOff:number
    // win more 开关
    SidebarOff:number
}


export class AdminInfo {
    /** 获取登录管理员详情 */
    static async GetAdminInfo(client, req: Empty): Promise<[AdminInfoResp, any]> {
        return await client.send("AdminInfo/GetAdminInfo", req)
    }

    /** 获取登录管理员详情 */
    static async GetTimeUtc(client, req: Empty): Promise<[AdminInfoResp, any]> {
        return await client.send("AdminInfo/GetTimeUtc", req)
    }

    /** 获取管理员列表 */
    static async AdminerList(client, req: PageListReq): Promise<[AdminListResp, any]> {
        return await client.send("AdminInfo/AdminerList", req)
    }

    /** 添加管理员 */
    static async AddAdminer(client, req: AddAdminReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddAdminer", req)
    }

    /** 添加商户管理员 */
    static async AddOperatorAdmin(client, req: AddOperatorAdminReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddOperatorAdmin", req)
    }

    /** 添加超级管理员 */
    static async AddAdmin(client, req: AddOperatorAdminReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddAdmin", req)
    }
    /** 添加商户 */
    static async AddMaintenance(client, req: AddMaintenance): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddOperatorV2", req)
    }
    /** 修改商户 */
    static async EditMaintenance(client, req: EditMaintenance): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdataOperatorV2", req)
    }
    /** 修改商户 */
    static async IncrementUpdataOperator(client, req: any): Promise<[Empty, any]> {
        return await client.send("AdminInfo/IncrementUpdataOperator", req)
    }
    /** 修改游戏内部配置 */
    static async UpdateGameUniversalConfig(client, req: UpdateGameUniversalConfigParam): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateGameUniversalConfig", req)
    }

    /** 添加商户用户 */
    static async AddOperatorUser(client, req: AddOperatorUserReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddOperatorUser", req)
    }

    //-----商户
    /** 添加商户 */
    static async AddOperator(client, req: AddOperatorReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddOperator", req)
    }

    /** 商户列表 */
    static async GetOperatorList(client, req: GetOperatorListReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetOperatorListV2", req)
    }
    /** 商户列表 */
    static async GetOperatorInfo(client, req: AddOperatorReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetOperatorById ", req)
    }
    /** 商户列表 */
    static async NewGetOperatorList(client, req: GetOperatorListReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/SearchOperator", req)
    }

    /** 商户列表 */
    static async GetOperatorChildList(client, req: GetOperatorListReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetLineOperator", req)
    }

    /** 玩家状态 */
    static async GetUpdatePlayerStatus(client, req: GetUpdatePlayerStatusReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdatePlayerStatus", req)
    }

    /** 获取商户详情 */
    static async getOperatorByAppId(client, req: AddOperatorReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/getOperatorByAppId", req)
    }

    /** 删除管理员 */
    static async DelAdminer(client, req: AdminPid): Promise<[Empty, any]> {
        return await client.send("AdminInfo/DelAdminer", req)
    }

    /** 删除管理员 */
    static async EditAdminer(client, req: AdminPid): Promise<[Empty, any]> {
        return await client.send("AdminInfo/EditAdmin", req)
    }

    /** 重置密码 */
    static async ResetPassword(client, req: AdminPid): Promise<[Empty, any]> {
        return await client.send("AdminInfo/ResetPassword", req)
    }

    /** 开关谷歌验证 */
    static async UpdateOpenGoole(client, req: UpdateOpenGooleReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateUserGoogleCode", req)
    }

    /** 冻结解冻管理员 */
    static async AdminerStatus(client, req: AdminerStatusReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateUserStatus", req)
    }

    /** 修改密码 */
    static async UpdatePasswd(client, req: UpdatePasswdReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdatePasswd", req)
    }

    /** 退出登录 */
    static async LoginOut(client, req: Empty): Promise<[Empty, any]> {
        return await client.send("AdminInfo/LoginOut", req)
    }

    /** 删除游戏 */
    static async DeleteGame(client, req: DeleteGameReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/DeleteGame", req)
    }

    /** 方案配置 */
    static async GetexcludeGameAll(client, req: GetExcludeGameAllReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/excludeGameAll", req)
    }

    static async GetExcludeGameAdd(client, req: GetExcludeGameAllReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/excludeGameAdd", req)
    }

    static async GetExcludeGameUp(client, req: GetExcludeGameAllReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/excludeGameUp", req)
    }

    static async GetExcludeGameDel(client, req: GetExcludeGameAllReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/excludeGameDel", req)
    }

    static async getServicesStatus(client, req: Empty): Promise<[Empty, any]> {
        return await client.send("mq/alerter/services/getStatus", req)
    }

    static async getSaveMsg4GameUpdateMsg(client, req: GetSaveMsg4GameUpdateMsgReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/saveMsg4GameUpdateMsg", req)
    }

    static async getSaveOperator4GameUpdateMsg(client, req: GetSaveOperator4GameUpdateMsgReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/saveOperator4GameUpdateMsg", req)
    }

    static async getConfig4GameUpdateMsg(client, req: {  }): Promise<[GetConfig4GameUpdateMsgResp, any]> {
        return await client.send("AdminInfo/getConfig4GameUpdateMsg", req)
    }

    static async getSendMsg4GameUpdate(client, req: GetSendMsg4GameUpdateReq): Promise<[GetConfig4GameUpdateMsgResp, any]> {
        return await client.send("AdminInfo/sendMsg4GameUpdate", req)
    }

    static async getResult4GameUpdateMsg(client, req: Empty): Promise<[GetResult4GameUpdateMsgResp, any]> {
        return await client.send("AdminInfo/Result4GameUpdateMsg", req)
    }

    static async getStatus4GameUpdateMsg(client, req: Empty): Promise<[GetStatus4GameUpdateMsgResp, any]> {
        return await client.send("AdminInfo/Status4GameUpdateMsg", req)
    }

}
