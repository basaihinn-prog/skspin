/* eslint-disable */
import type { Empty } from "../empty";


export interface FileContent {
  FileName?: string | undefined;
  Content?: string | undefined;
}

export interface DBConfigFileLog {
  /** @gotags:bson:"_id" */
  Id?:
    | string
    | undefined;
  /** 操作人 */
  User?: number | undefined;
  FileName?: string | undefined;
  Content?:
    | string
    | undefined;
  /** 操作时间 */
  CreateAt?: string | undefined;
}

export interface FileHistoryResp {
  List?: DBConfigFileLog[] | undefined;
}

export interface FileHistoryReq {
  /** 最近多少条 */
  Count?:
    | number
    | undefined;
  /** 文件名 */
  FileName?: string | undefined;
}

export class AdminConfigFile {
  /** 保存配置 */
  static async SaveConfig(client, req : FileContent) : Promise<[Empty,any]> {
    return await client.send("AdminConfigFile/SaveConfig", req)
  }
  /** 读取配置 Content 不传 */
  static async LoadConfig(client, req : FileContent) : Promise<[FileContent,any]> {
    return await client.send("AdminConfigFile/LoadConfig", req)
  }
  /** 查看某个文件的历史记录 最近n条 */
  static async FileHistory(client, req : FileHistoryReq) : Promise<[FileHistoryResp,any]> {
    return await client.send("AdminConfigFile/FileHistory", req)
  }
}