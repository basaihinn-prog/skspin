/* eslint-disable */
import type { Empty } from "../empty";
import type { AdminStatus, DBAdminer } from "./auth";


/** 权限组表 */
export interface DBAuth {
  /** @gotags:bson:"_id" */
  ID?:
    | number
    | undefined;
  /** 菜单id */
  MenuIds?:
    | string
    | undefined;
  /** 权限组名称 */
  Name?:
    | string
    | undefined;
  /** 备注 */
  Remark?:
    | string
    | undefined;
  /** 创建时间 */
  CreateTime?: string | undefined;
}

export interface AdminListResp {
  Count?: number | undefined;
  List?: DBAdminer[] | undefined;
}

export interface AddAdminReq {
  Username?: string | undefined;
  GroupId?: number | undefined;
}

export interface AdminPid {
  Pid?: number | undefined;
}

export interface UpdateOpenGooleReq {
  Pid?: number | undefined;
  IsOpenGoole?: boolean | undefined;
}

export interface AdminerStatusReq {
  Pid?: number | undefined;
  Status?: AdminStatus | undefined;
}

/** 菜单表 */
export interface DBMenu {
  /** @gotags:bson:"_id" */
  ID?:
    | number
    | undefined;
  /** 菜单标题 */
  Title?:
    | string
    | undefined;
  /** 菜单pid 0为顶级pid */
  Pid?:
    | number
    | undefined;
  /** 节点 */
  Url?:
    | string
    | undefined;
  /** 图标icon */
  Icon?:
    | string
    | undefined;
  /** 排序 */
  Sort?: number | undefined;
}

export interface AdminInfoResp {
  ID?: number | undefined;
  Username?: string | undefined;
  Gid?: number | undefined;
  MenuList?: DBMenu[] | undefined;
}

export interface UpdatePasswdReq {
  OldPasswd?: string | undefined;
  NewPasswd?: string | undefined;
  ConfirmPasswd?: string | undefined;
}

export interface PageListReq {
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
  Query?: string | undefined;
  Sort?: string | undefined;
}

export class AdminInfo {
  /** 获取登录管理员详情 */
  static async GetAdminInfo(client, req : Empty) : Promise<[AdminInfoResp,any]> {
    return await client.send("AdminInfo/GetAdminInfo", req)
  }
  /** 获取管理员列表 */
  static async AdminerList(client, req : PageListReq) : Promise<[AdminListResp,any]> {
    return await client.send("AdminInfo/AdminerList", req)
  }
  /** 添加管理员 */
  static async AddAdminer(client, req : AddAdminReq) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/AddAdminer", req)
  }
  /** 删除管理员 */
  static async DelAdminer(client, req : AdminPid) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/DelAdminer", req)
  }
  /** 开关谷歌验证 */
  static async UpdateOpenGoole(client, req : UpdateOpenGooleReq) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/UpdateOpenGoole", req)
  }
  /** 冻结解冻管理员 */
  static async AdminerStatus(client, req : AdminerStatusReq) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/AdminerStatus", req)
  }
  /** 修改密码 */
  static async UpdatePasswd(client, req : UpdatePasswdReq) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/UpdatePasswd", req)
  }
  /** 退出登录 */
  static async LoginOut(client, req : Empty) : Promise<[Empty,any]> {
    return await client.send("AdminInfo/LoginOut", req)
  }
}