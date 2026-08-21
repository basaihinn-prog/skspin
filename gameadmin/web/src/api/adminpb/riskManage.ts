import {NotifyList} from "@/api/adminpb/notify";
import {Empty} from "@/api/empty";

const ADD_RISK_RULE = '/AdminInfo/AddRickRule'
const GET_RISK_RULE = '/AdminInfo/GetRickRule'


export interface RiskRuleRequest{
    AppID: string
    RickRule?: string
    Type?: string
}


export interface RuleAggregate {
    // 规则集合名称
    RuleAggregateName: string;
    // 有效时间
    EffectTime?: string;
    // 失效时间
    InvalidTime?: string;
    TimeRange?: string[];
    // 状态
    Status: number;
    // 规格集合
    Rules: string;

}


export interface RiskRuleResponse{
    AppID      :string;
    ReturnRate :any;
    OriginReturnRate :string;
    Transfer   :any;
    OriginTransfer   :string;
}



export class Rule {
    static async AddRiskRule(client, req:RiskRuleRequest): Promise<[Empty, any]> {
        return await client.send(ADD_RISK_RULE, req)
    }
    static async GetRiskRule(client, req:RiskRuleRequest): Promise<[RiskRuleResponse, any]> {
        return await client.send(GET_RISK_RULE, req)
    }
}
