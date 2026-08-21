/* eslint-disable */
import type { Empty } from "../empty";
import type { DBAuth, DBMenu, PageListReq } from "./info";


export interface GroupListResp {
  Count?: number | undefined;
  List?: DBAuth[] | undefined;
}

export interface AddGroupReq {
  Name?: string | undefined;
  Remark?: string | undefined;
}

export interface GroupIdReq {
  Gid?: number | undefined;
}

export interface EditPowerReq {
  PowerList?: number[] | undefined;
  Gid?: number | undefined;
}

export interface PowerListResp {
  MenuList?: DBMenu[] | undefined;
  Gid?: number | undefined;
  CheckedArray?: number[] | undefined;
}

export class AdminGroup {
  /** 获取权限组列表 */
  static async GroupList(client, req : PageListReq) : Promise<[GroupListResp,any]> {
    return await client.send("AdminGroup/GroupList", req)
  }
  /** 添加权限组 */
  static async AddGroup(client, req : AddGroupReq) : Promise<[Empty,any]> {
    return await client.send("AdminGroup/AddGroup", req)
  }
  /** 根据权限ID查看权限详情 */
  static async GetPowerDetailsById(client, req : GroupIdReq) : Promise<[PowerListResp,any]> {
    return await client.send("AdminGroup/GetPowerDetailsById", req)
  }
  /** 修改权限组菜单 */
  static async EditPowerDetails(client, req : EditPowerReq) : Promise<[Empty,any]> {
    return await client.send("AdminGroup/EditPowerDetails", req)
  }
}