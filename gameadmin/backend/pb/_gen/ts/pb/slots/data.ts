/* eslint-disable */


export interface ThreeMonkeyRunTimeData {
  Channel?: string | undefined;
  Date?: string | undefined;
  Expense?: number | undefined;
  Income?: number | undefined;
  NetProfit?: number | undefined;
  SelfPoolReward?: number | undefined;
  Hit?: number | undefined;
  AwardHit?: number | undefined;
  SmallGameHit?: number | undefined;
  SmallGameIncome?: number | undefined;
  FreeGameHit?: number | undefined;
  FreeGameIncome?: number | undefined;
  EnterCount?: number | undefined;
  WinPlrCount?: number | undefined;
  PlrCount?: number | undefined;
}

export interface ThreeMonkeyRunTimeResp {
  Arr?: ThreeMonkeyRunTimeData[] | undefined;
  Pools?: string[] | undefined;
  GlobalN?: number | undefined;
  TotalExpense?: number | undefined;
  TotalIncome?: number | undefined;
}

export interface SlotsStatsBuyGame {
}

export interface SlotsStatsBuyGame_Doc {
  /** @gotags: bson:"_id" */
  ID?: string | undefined;
  Day?: string | undefined;
  Game?: string | undefined;
  Bet?: number | undefined;
  Win?: number | undefined;
}

export interface SlotsStatsBuyGame_QueryReq {
  Days?: string[] | undefined;
  Game?: string | undefined;
}

export interface SlotsStatsBuyGame_QueryResp {
  List?: SlotsStatsBuyGame_Doc[] | undefined;
}