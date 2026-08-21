export const OperatorRTP = ["", "10-97", "93-300", "10-300"]
export const ServiceAreaList = [
    {label:"美洲地区",value:"US", queryLocal:import.meta.env.VITE_US_REQUEST_URI},
    {label:"亚太地区",value:"AS", queryLocal:import.meta.env.VITE_AS_REQUEST_URI}
]

export const TimeZone = [
    {label:"GMT+0",value:0},
    {label:"GMT+1",value:3600000},
    {label:"GMT+2",value:7200000},
    {label:"GMT+3",value:10800000},
    {label:"GMT+4",value:14400000},
    {label:"GMT+5",value:18000000},
    {label:"GMT+6",value:21600000},
    {label:"GMT+7",value:25200000},
    {label:"GMT+8",value:28800000},
    {label:"GMT+9",value:32400000},
    {label:"GMT+10",value:36000000},
    {label:"GMT+11",value:39600000},
    {label:"GMT+12",value:43200000},
    {label:"GMT-1",value:-3600000},
    {label:"GMT-2",value:-7200000},
    {label:"GMT-3",value:-10800000},
    {label:"GMT-4",value:-14400000},
    {label:"GMT-5",value:-18000000},
    {label:"GMT-6",value:-21600000},
    {label:"GMT-7",value:-25200000},
    {label:"GMT-8",value:-28800000},
    {label:"GMT-9",value:-32400000},
    {label:"GMT-10",value:-36000000},
    {label:"GMT-11",value:-39600000},
    {label:"GMT-12",value:-43200000},
]




export const GameType = [
    {label: "心跳型", value: 1},
    {label: "波动型", value: 2},
    {label: "仿正版", value: 3},
    {label: "混合型", value: 4},
    {label: "稳定型", value: 5},
    {label: "高中奖率", value: 6},
    {label: "超高中奖率", value: 7},
]

export const BetMultiples = [0.2,0.5,1,2,3,5,10,20,50,100,200,500,1000]

export const GameClass: JsMap<number, string> = {
    0: 'SLOT',
    1: 'MINI',

};
