/** 添加多语言 */
import {DBAuth} from "@/api/adminpb/info";
import {Empty} from "@/api/empty";
const AddNotification = "AdminInfo/addNotification"
const getEffectiveNotify = "AdminInfo/getEffectiveNotify"
const getNotification = "AdminInfo/getNotification"
const editNotification = "AdminInfo/editNotification"
const deleteNotification = "AdminInfo/deleteNotification"


export interface NotifyParams {
    notifyTitle: string,
    notifyType: string,
    notifyStatus: string,
    notifyCreateTime: string,
    onlineLimitTime: string,
    lowLimitTime: string,
    Sort: number,
    languageContext: string
}


export interface NotifyList {
    Count?: number;
    List?: DBAuth[];
}


export class Notify {


    static async GetNotify(client, req: NotifyParams): Promise<[NotifyList, any]> {
        return await client.send(getNotification, req)
    }

    static async GetEffectiveNotify(client, req: NotifyParams): Promise<[NotifyList, any]> {
        return await client.send(getEffectiveNotify, req)
    }

    static async AddNotify(client, req: NotifyParams): Promise<[Empty, any]> {
        return await client.send(AddNotification, req)
    }
    static async EditNotify(client, req: NotifyParams): Promise<[Empty, any]> {
        return await client.send(editNotification, req)
    }
    static async DelNotify(client, req: NotifyParams): Promise<[Empty, any]> {
        return await client.send(deleteNotification, req)
    }

}
