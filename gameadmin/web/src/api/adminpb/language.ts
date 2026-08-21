import {Empty} from "@/api/empty";
import {FileContent, FileHistoryReq, FileHistoryResp} from "@/api/adminpb/json";
import {CommPageRequest} from "@/api/comm";

const GET_LANG = "AdminInfo/GetConfig"
const UPLOAD_LANG = "AdminInfo/uploadConfig"
const EDIT_LANG = "AdminInfo/EditConfig"
const SELECT_LANG = "AdminInfo/SelectLang"

export interface LanguageReq extends CommPageRequest{
    FileName: string
    LanguageConfig: Object[]
    Context:Object[]
}

export class Language {

    static async EditLanguageConfig(client, req: LanguageReq): Promise<[LanguageReq,any]> {
        return client.send(EDIT_LANG, req)
    }
    static async SelectConfig(client, req: LanguageReq): Promise<[LanguageReq,any]> {
        return client.send(GET_LANG, req)
    }
    static async UploadConfig(client, req: LanguageReq): Promise<[LanguageReq,any]> {
        return client.send(UPLOAD_LANG, req)
    }

}
