
// 获取预警接口
import {CommPageRequest} from "@/api/comm";

const GET_RETURN_RATE_HISTORY = "AdminInfo/getAlertDataHistory"
const SET_RETURN_RATE_HISTORY = "AdminInfo/setAlertDataHistory"


export interface AlarmRequest extends CommPageRequest{
    Pid?: number
    GameId?: string
    AppId?: string
    Manufacturer?: string
    AlarmTimeStart?: number
    AlarmTimeEnd?: number
    Type: string
}


export interface EditRequest{
    Type?: string
    Id?: string
    Range?: string
}



export interface AlarmListResponse {
    List: AlarmData[]
    Count: number
}


export interface AlarmData {
    Msg: string
    Pid: number
    UserId: string
    WinMoney: number
    ReadStatus: number
    Id: string,
    GameId:string,
    OrderId: string,
    TotalWinLoss: number,
    TotalBet: number,
    WinRate: number,
    Balance: string,
    AppId:  string,
    Amount:  string,
    CreateTime: string

}


export interface AlarmItem extends AlarmData{

    OrderId: string
    CreateTime: string
    ProfitAndLoss: string
    BetMount: string
    LogInfo: string
    TotalRate: string
    Currency: string
    MerchantId: string
}


export class Alarm {
    static async GetAlarmHistory(client, req: AlarmRequest): Promise<[AlarmListResponse, any]> {
        return await client.send(GET_RETURN_RATE_HISTORY, req)
    }
    static async SetAlarmHistory(client, req: EditRequest): Promise<[AlarmListResponse, any]> {
        return await client.send(SET_RETURN_RATE_HISTORY, req)
    }
}
