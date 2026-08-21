/* eslint-disable */
import type { Empty } from "../empty";


export interface FileContent {
  FileName?: string;
  Content?: string;
}

export interface DBConfigFileLog {
  /** @gotags:bson:"_id" */
  Id?: string;
  /** 操作人 */
  User?: number;
  FileName?: string;
  Content?: string;
  /** 操作时间 */
  CreateAt?: string;
}

export interface FileHistoryResp {
  List?: DBConfigFileLog[];
}

export interface FileHistoryReq {
  /** 最近多少条 */
  Count?: number;
  /** 文件名 */
  FileName?: string;
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