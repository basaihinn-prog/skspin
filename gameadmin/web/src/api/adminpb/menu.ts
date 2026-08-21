/* eslint-disable */
import type { Empty } from "../empty";
import type { DBMenu } from "./info";


export interface MenuListResp {
  MenuList?: DBMenu[];
}

export interface DelMenuReq {
  ID?: number;
}

export class AdminMenu {
  /** 获取菜单列表 */
  static async MenuList(client, req : Empty) : Promise<[MenuListResp,any]> {
    return await client.send("AdminMenu/MenuList", req)
  }
  /** 添加菜单 */
  static async AddMenu(client, req : DBMenu) : Promise<[Empty,any]> {
    return await client.send("AdminMenu/AddMenu", req)
  }
  /** 编辑菜单 */
  static async UpdateMenu(client, req : DBMenu) : Promise<[Empty,any]> {
    return await client.send("AdminMenu/UpdateMenu", req)
  }
  /** 删除菜单 */
  static async DelMenu(client, req : DelMenuReq) : Promise<[Empty,any]> {
    return await client.send("AdminMenu/DelMenu", req)
  }
}