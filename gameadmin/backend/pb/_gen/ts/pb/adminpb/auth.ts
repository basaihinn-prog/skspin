/* eslint-disable */


export enum AdminStatus {
  /** Normal - 正常 */
  Normal = 0,
  /** Frozen - 冻结 */
  Frozen = 1,
}

export enum IsDelete {
  /** NoDelete - 正常 */
  NoDelete = 0,
  /** YesDelete - 删除 */
  YesDelete = 1,
}

export interface AdminLoginReq {
  /** 用户名 */
  Username?:
    | string
    | undefined;
  /** 密码 */
  Password?:
    | string
    | undefined;
  /** 谷歌验证码 */
  GooleCode?: string | undefined;
}

export interface AdminLoginResp {
  Token?: string | undefined;
  ExpireAt?: string | undefined;
}

export interface DBAdminer {
  /** @gotags: bson:"_id" */
  ID?:
    | number
    | undefined;
  /** 用户名（昵称） */
  Username?:
    | string
    | undefined;
  /** @gotags: json:"-" */
  Password?:
    | string
    | undefined;
  /** 谷歌验证码 */
  GoogleCode?:
    | string
    | undefined;
  /** 谷歌二维码 */
  Qrcode?:
    | string
    | undefined;
  /** 头像 */
  Avatar?:
    | string
    | undefined;
  /** 是否打开谷歌验证码开关 */
  IsOpenGoole?:
    | boolean
    | undefined;
  /** 管理员状态 */
  Status?:
    | AdminStatus
    | undefined;
  /** 是否删除 */
  IsDelete?:
    | IsDelete
    | undefined;
  /** 注册时间 */
  CreateAt?:
    | string
    | undefined;
  /** 登录时间 */
  LoginAt?:
    | string
    | undefined;
  /** 权限组ID */
  GroupId?:
    | number
    | undefined;
  /** 最后登陆的token，用于单点登录 */
  Token?:
    | string
    | undefined;
  /** Token有效期 */
  TokenExpireAt?: string | undefined;
}

export class AdminAuth {
  static async Login(client, req : AdminLoginReq) : Promise<[AdminLoginResp,any]> {
    return await client.send("AdminAuth/Login", req)
  }
}