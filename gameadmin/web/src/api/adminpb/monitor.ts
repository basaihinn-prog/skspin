import {Empty} from "@/api/empty";

const SET_GAME_MONITOR = "/AdminInfo/SetGameMonitor"

export interface SetGameMonitorParams {
    AppID                      : string
    MoniterType                : number
    IsMoniter                  : number
    MoniterNewbieNum           : number
    MoniterRTPErrorValue       : number
    MoniterNumCycle            : number
    MoniterAddRTPRangeValue    : MoniterRtpRangeValue[]
    MoniterReduceRTPRangeValue : MoniterRtpRangeValue[]
}

export interface RTPProtectInfo {

    //游戏监控
    IsMoniter: number
    // 监控新手数量
    MoniterNewbieNum: number
    // 监控RTP错误值
    MoniterRTPErrorValue: number
    // 监控周期
    MoniterNumCycle: number
    // 监控增加RTP范围值
    MoniterAddRTPRangeValue: MoniterRtpRangeValue[]
    // 监控减少RTP范围值
    MoniterReduceRTPRangeValue: MoniterRtpRangeValue[]
}

export interface MoniterRtpRangeValue {

    RangeMinValue  :number
    RangeMaxValue  :number
    NewbieValue    :number
    NotNewbieValue :number
}

export class Monitor{
    static async SetGameMonitor(client, req: SetGameMonitorParams): Promise<[Empty, any]> {
        return await client.send(SET_GAME_MONITOR, req)
    }
}
