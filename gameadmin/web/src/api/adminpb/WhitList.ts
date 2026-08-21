import {App} from "@/api/gamepb/admin";
import {Empty} from "@/api/empty";
import {CommPageRequest} from "@/api/comm";

export  interface WhitResponse {
    List?: WhitDataResponse[];
    Count:number
}

export interface WhitData {
    UserIp: string;
    Remark: string;
}


export interface WhitDataResponse extends WhitData {
    OperatorTime: string;
    OperatorName: string;
}



export class White {

    /** 添加白名单 */
    static async AddWhitUser(client, req: WhitData): Promise<[Empty, any]> {
        return await client.send("AdminInfo/add_whitList", req)
    }

    /** 查询白名单 */
    static async GetWhitUserList(client, req: CommPageRequest): Promise<[WhitResponse, any]> {
        return await client.send("AdminInfo/get_whitList", req)
    }

    /** 删除白名单 */
    static async DeleteWhitUser(client, req: any): Promise<[Empty, any]> {
        return await client.send("AdminInfo/del_system_maintenance", req)
    }


}
