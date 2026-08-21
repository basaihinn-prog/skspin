import {Empty} from "@/api/empty";


export interface CommPageRequest {
    PageSize: number,
    Page: number
}
export interface UploadResponse {
    Url?: string;
}

export interface Currency{

    CurrencyName: string
    CurrencyCode: string
    CurrencySymbol: string

}

export class Upload {
    /** 金币操作 */
    static async UploadFile(client, req : FormData) : Promise<[UploadResponse,any]> {
        return await client.send("comm/saveUpload", req)
    }

}


export class Comm {
    /** 币种查询 */
    static async GetCurrency(client, req : FormData) : Promise<[Currency[],any]> {
        return await client.send("AdminInfo/GetCurrency", req)
    }
}
