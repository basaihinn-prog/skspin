/* eslint-disable */


export interface ThreeMonkeyRunTimeData {
  Channel?: string;
  Date?: string;
  Expense?: number;
  Income?: number;
  NetProfit?: number;
  SelfPoolReward?: number;
  Hit?: number;
  AwardHit?: number;
  SmallGameHit?: number;
  SmallGameIncome?: number;
  FreeGameHit?: number;
  FreeGameIncome?: number;
  EnterCount?: number;
  WinPlrCount?: number;
  PlrCount?: number;
}

export interface ThreeMonkeyRunTimeResp {
  Arr?: ThreeMonkeyRunTimeData[];
  Pools?: string[];
  GlobalN?: number;
  TotalExpense?: number;
  TotalIncome?: number;
}

export interface SlotsStatsBuyGame {
}

export interface SlotsStatsBuyGame_Doc {
  /** @gotags: bson:"_id" */
  ID?: string;
  Day?: string;
  Game?: string;
  Bet?: number;
  Win?: number;
}

export interface SlotsStatsBuyGame_QueryReq {
  Days?: string[];
  Game?: string;
}

export interface SlotsStatsBuyGame_QueryResp {
  List?: SlotsStatsBuyGame_Doc[];
}