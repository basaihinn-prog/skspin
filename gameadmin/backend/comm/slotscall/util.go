package slotscall

import (
	"fmt"
	"game/pb/_gen/pb/slots"
)

func ThreeMonkeyRunTimeData_Add(this *slots.ThreeMonkeyRunTimeData, delta *slots.ThreeMonkeyRunTimeData) {
	if delta == nil || this == nil {
		return
	}
	if this.Channel == "" {
		this.Channel = delta.Channel
	} else {
		this.Channel = fmt.Sprintf("%s,%s", this.Channel, delta.Channel)
	}
	this.Date = delta.Date
	this.Expense += delta.Expense
	this.Income += delta.Income
	this.NetProfit += delta.NetProfit
	this.SelfPoolReward += delta.SelfPoolReward
	this.Hit += delta.Hit
	this.AwardHit += delta.AwardHit
	this.SmallGameHit += delta.SmallGameHit
	this.SmallGameIncome += delta.SmallGameIncome
	this.FreeGameHit += delta.FreeGameHit
	this.FreeGameIncome += delta.FreeGameIncome
	this.EnterCount += delta.EnterCount
	this.WinPlrCount += delta.WinPlrCount
	this.PlrCount += delta.PlrCount
}

func ThreeMonkeyRunTimeResp_Add(this *slots.ThreeMonkeyRunTimeResp, delta *slots.ThreeMonkeyRunTimeResp) {
	if delta == nil {
		return
	}

	if len(this.Arr) != len(delta.Arr) {
		return
	}
	for i, v := range delta.Arr {
		if this.Arr[i] == nil {
			this.Arr[i] = &slots.ThreeMonkeyRunTimeData{}
		}
		ThreeMonkeyRunTimeData_Add(this.Arr[i], v)
	}
	this.TotalIncome += delta.TotalIncome
	this.TotalExpense += delta.TotalExpense
}
