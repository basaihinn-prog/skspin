/* eslint-disable */


export enum SearchType {
  Day = 0,
  Mon = 1,
}

export interface PlayerAnalysis {
  /** @gotags: bson:"_id" */
  ID?: string | undefined;
  Pid?:
    | number
    | undefined;
  /**  */
  Bet?:
    | number
    | undefined;
  /**  */
  Win?: number | undefined;
}

export interface AppAnalysis {
  /** @gotags: bson:"_id" */
  ID?: string | undefined;
  AppID?: string | undefined;
  Bet?: number | undefined;
  Win?: number | undefined;
}

export interface GameAnalysis {
  /** @gotags: bson:"_id" */
  ID?: string | undefined;
  GameID?: string | undefined;
  Bet?: number | undefined;
  Win?: number | undefined;
}

export interface PlayerAnalysisReq {
  Pid?: number | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
}

export interface PlayerAnalysisResp {
  Count?: number | undefined;
  List?: PlayerAnalysis[] | undefined;
}

export interface AppAnalysisReq {
  AppID?: string | undefined;
  Search?: SearchType | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
}

export interface AppAnalysisResp {
  Count?: number | undefined;
  List?: AppAnalysis[] | undefined;
}

export interface GameAnalysisReq {
  GameID?: string | undefined;
  Search?: SearchType | undefined;
  PageIndex?: number | undefined;
  PageSize?: number | undefined;
}

export interface GameAnalysisResp {
  Count?: number | undefined;
  List?: GameAnalysis[] | undefined;
}

/** @ts: prefix(gamecenter) */
export class AdminAnalysis {
  static async Player(client, req : PlayerAnalysisReq) : Promise<[PlayerAnalysisResp,any]> {
    return await client.send("gamecenter/gc.AdminAnalysis/Player", req)
  }
  static async App(client, req : AppAnalysisReq) : Promise<[AppAnalysisResp,any]> {
    return await client.send("gamecenter/gc.AdminAnalysis/App", req)
  }
  static async Game(client, req : GameAnalysisReq) : Promise<[GameAnalysisResp,any]> {
    return await client.send("gamecenter/gc.AdminAnalysis/Game", req)
  }
}