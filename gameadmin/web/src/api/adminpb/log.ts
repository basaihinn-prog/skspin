import {CommPageRequest, Currency} from "@/api/comm";
import {Empty} from "@/api/empty";

const ADD_SYSTEM_LOG = "/AdminInfo/AddLog"
const GET_LOG = "/AdminInfo/GetLog"


// 系统日志列表查询
export interface SystemLogListParams extends CommPageRequest{
    OperateMod:   number,
    StartTime:    number,
    EndTime:      number,
    OperateType:  number,
    OperatorName: string,
}

export interface SystemLogListResponse {
    List :SystemLog[]
    All  :number
}


export interface SystemLog {
    ID             : string
    OperatorID     : number                 //操作人id
    OperatorName   : string                 //操作人名称
    OperatorType   : number                 //操作人类型
    OperateType    : number                 //操作类型
    OperateContent : string                 //操作内容
    OperateTime    : string                 //操作时间
    IpAdress       : string                 //操作ip
    OperateMod     : number                 //操作模块
    AppID          : string                 //运营商
}





// 系统日志列表
export interface AddSystemLogParam {
    OperateMod:     number,
    OperateType:    number,
    IpAdress:       string,
    OperateContent: string,
}

// 日志模块枚举值
export enum LOG_PRIMARY_MODULE {
    LOGIN = 1,
    OPERATOR_LIST,
    ADD_OPERATOR,
    PLAYER_INFO,
    OTHER_RTP,
    GAME_CONTROL,
    PLAYER_CONTROL,
    GAME_LIST,
    ABILITY_MENU,
    RULE,
}


// 日志操作枚举值
export enum LOG_OPERATOR_TYPE {

    LOGIN = 1,
    EXIT,
    RESET,
    CLOSE,
    ADD,
    REMOVE,
    EDIT,
    SELECT,
    SELECT_INFO,


    BALANCE_EDIT = 999        // 特殊操作:修改余额

}



export class Log {
    /** 币种查询 */
    static async AddSystemLog(client, req : AddSystemLogParam) : Promise<[Empty,any]> {
        return await client.send(ADD_SYSTEM_LOG, req)
    }
    /** 币种查询 */
    static async GetSystemLog(client, req : SystemLogListParams) : Promise<[SystemLogListResponse[],any]> {
        return await client.send(GET_LOG, req)
    }


}
