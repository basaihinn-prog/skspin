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
  Username?: string;
  /** 密码 */
  Password?: string;
  /** 谷歌验证码 */
  GooleCode?: string;
}

export interface AdminLoginResp {
  Token?: string;
  ExpireAt?: string;
}

export interface DBAdminer {
  /** @gotags: bson:"_id" */
  ID?: number;
  /** 用户名（昵称） */
  Username?: string;
  /** @gotags: json:"-" */
  Password?: string;
  /** 谷歌验证码 */
  GoogleCode?: string;
  /** 谷歌二维码 */
  Qrcode?: string;
  /** 头像 */
  Avatar?: string;
  /** 是否打开谷歌验证码开关 */
  IsOpenGoole?: boolean;
  /** 管理员状态 */
  Status?: AdminStatus;
  /** 是否删除 */
  IsDelete?: IsDelete;
  /** 注册时间 */
  CreateAt?: string;
  /** 登录时间 */
  LoginAt?: string;
  /** 权限组ID */
  GroupId?: number;
  /** 最后登陆的token，用于单点登录 */
  Token?: string;
  /** Token有效期 */
  TokenExpireAt?: string;
}

export class AdminAuth {
  static async Login(client, req : AdminLoginReq) : Promise<[AdminLoginResp,any]> {
    return await client.send("AdminAuth/Login", req)
  }
}