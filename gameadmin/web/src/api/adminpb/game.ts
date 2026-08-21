import {PageListReq} from "@/api/adminpb/info";
import {GroupListResp} from "@/api/adminpb/group";
import {Empty} from "@/api/empty";

export interface GameTypeListRep {
    List: GameTypeRep[];
    Count: number;
}

export interface GameTypeRep {
    GameTypeName: string;
    GameNum: number;
    GameEditStatus: boolean;
}

export interface GameTypeReq {
    GameTypeName: string;
}



export class GameType {


    static async GameTypeList(client, req: PageListReq): Promise<[GameTypeListRep, any]> {
        return await client.send("AdminInfo/GetPermissionList", req)
    }

    static async AddGameType(client, req: GameTypeReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetPermissionList", req)
    }

    static async EditGameType(client, req: GameTypeReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetPermissionList", req)
    }

    static async DeleteGameType(client, req: GameTypeReq): Promise<[Empty, any]> {
        return await client.send("AdminInfo/GetPermissionList", req)
    }

}

