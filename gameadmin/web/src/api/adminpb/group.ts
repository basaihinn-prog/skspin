/* eslint-disable */
import type {Empty} from "../empty";
import type {DBAuth, DBMenu, PageListReq} from "./info";


export interface GroupListResp {
    Count?: number;
    List?: DBAuth[];
}

export interface AddLangResp {
    ErrLangs:[]
}

export interface AddGroupReq {
    PermissionId?: number
    Name?: string
    Remark?: string
    MenuIds?: number[]
}

export interface GroupIdReq {
    PermissionId?: number
    Name?: string
    Remark?: string
    MenuIds?: number[]
}

export interface EditPowerReq {
    PowerList?: number[];
    Gid?: number;
}

export interface PowerListResp {
    MenuList?: DBMenu[];
    Gid?: number;
    CheckedArray?: number[];
}

export interface EditOperatorReq {
    MenuIds?: DBMenu[];
    OperatorId?: number;
}

export interface DownloadOperatorDataReq {
    OperatorId?: number;
}

export interface GetAppGameInfoReq {
    Operator?: number;
    Date?: string;
}

export interface UpdateOperatorStatusReq{
    OperatorId:number;
    Status:number
}

export interface GetLangReq {
    PageIndex?: number;
    PageSize?: number;
}
export interface AddLangReq {
    Langs:[]
}
export interface LangDeleteReq {
    ZH:string
}

export interface MonthBillList {
    List: MonthBill[]
}
export interface MonthBill {
    AppId: string
    Bet: number
    Game: string
    Time: string
    Win: number
    Profit:number
    PlantRate:number
    Receivable:number
    ChangeReceivable:number
    TurnoverPay:number
    CooperationType:number
    exchangeRate?:number
    OperatorType?:number
    ChangePlantRateReceivable:number
    ChangeTurnoverPayReceivable:number
    children: MonthBill[]
}


export class AdminGroup {
    /** 获取权限组列表 */
    static async GroupList(client, req: PageListReq): Promise<[GroupListResp, any]> {
        return await client.send("AdminInfo/GetPermissionList", req)
    }
    static async GroupListAll(client, req: PageListReq): Promise<[GroupListResp, any]> {
        return await client.send("AdminInfo/GetPermissionListAll", req)
    }

    /** 添加多语言 */
    static async UploadLang(client, req: FormData): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddLangToDB", req)
    }
    /** 添加多语言 */
    static async GetLang(client, req: GetLangReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/LangGet", req)
    }
    /** 添加多语言 */
    static async GetCurrency (client, req: any): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetCurrency", req)
    }

    /** 添加多语言 */
    static async AddLang(client, req: AddLangReq): Promise<[AddLangResp, any]> {
        return await client.send("AdminInfo/LangAdd", req)
    }

    /** 获取某行多语言 */
    static async LangUpdate(client, req: AddLangReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/LangUpdate", req)
    }
    /** 获取某行多语言 */
    static async LangDelete(client, req: LangDeleteReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/LangDelete", req)
    }

    /** 添加权限组 */
    static async AddGroup(client, req: AddGroupReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/AddPermission", req)
    }

    /** 根据权限ID查看权限详情 */
    static async GetPowerDetailsById(client, req: GroupIdReq): Promise<[PowerListResp, any]> {
        return await client.send("AdminInfo/UpdatePermission", req)
    }

    /** 修改权限组菜单 */
    static async EditPowerDetails(client, req: EditPowerReq): Promise<[Empty, any]> {
        return await client.send("AdminGroup/EditPowerDetails", req)
    }

    /** 修改商户 */
    static async EditOperator(client, req: EditOperatorReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateOperatorPermission", req)
    }

    static async DownloadOperatorData(client, req: DownloadOperatorDataReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/DownloadOperatorData", req)
    }

    static async GetAppGameInfo(client, req: GetAppGameInfoReq): Promise<[MonthBillList, any]> {
        return await client.send("AdminInfo/GetAppGameInfo", req)
    }

    static async UpdateOperatorStatus(client, req: UpdateOperatorStatusReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/UpdateOperatorStatus", req)
    }
}
