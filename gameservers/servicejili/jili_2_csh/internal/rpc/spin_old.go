package rpc

import (
	"serve/comm/db"
	"serve/comm/define"
	"serve/comm/jwtutil"
	"serve/comm/slotsmongo"
	"serve/comm/ut"
	"serve/servicejili/jili_2_csh/internal/message"
	"serve/servicejili/jili_2_csh/internal/models"

	"google.golang.org/protobuf/proto"
)

func spin_old(ps define.JiliParams, ret *[]byte) (err error) {
	var gameReq message.Csh_GameReqData
	err = proto.Unmarshal(ps.Data, &gameReq)
	if err != nil {
		return
	}
	token := gameReq.GetToken()

	pid, err := jwtutil.ParseToken(token)
	if err != nil {
		// err = define.NewErrCode("Invalid player session", 1302)
		return
	}
	gold, err := slotsmongo.GetBalance(pid)
	if err != nil {
		return err
	}

	var spinReq message.Csh_SpinReq
	err = proto.Unmarshal(gameReq.Encode, &spinReq)
	if err != nil {
		return
	}

	return db.CallWithPlayer(pid, func(plr *models.Player) error {
		var spinAll = message.Csh_SpinAllData{
			Info: []*message.Csh_SpinAck{
				{
					TotalWin:  proto.Float64(1.9),
					MGWin:     proto.Float64(1.9),
					NowMoney:  proto.Float64(ut.Gold2Money(gold)),
					ShowIndex: proto.String("17163-756975-05404002"),
					// AwardTypeFlag: proto.Int32(1),
					AckVec: []*message.Csh_SinglePlate{
						{
							AwardTypeFlag: proto.Int32(1),
							Coin:          proto.Float64(ut.Gold2Money(gold)),
							ExtraOdds:     proto.Int32(1),
							LogIndex:      proto.String("1716375697505404002"),
							PlateVec:      []*message.Csh_Column{},
							RemainRound:   proto.Int32(8),
							AwardVec: []*message.Csh_AwardDetail{
								{
									AwardMoney: proto.Float64(0.3),
									AwardSet: []*message.Csh_AwardData{
										{
											AwardMoney:  proto.Float64(0.1),
											AwardSymbol: proto.Int32(6),
											ContainWild: proto.Bool(true),
											Lines:       proto.Int32(1),
											SymbolCount: proto.Int32(3),

											AwardGridVec: []*message.Csh_Grid{
												{
													Row: proto.Int32(1),
												},
												{
													Column: proto.Int32(1),
													Row:    proto.Int32(1),
												},
												{
													Column: proto.Int32(2),
												},
											},
										},
										{
											AwardMoney:  proto.Float64(0.2),
											AwardSymbol: proto.Int32(7),
											ContainWild: proto.Bool(true),
											Lines:       proto.Int32(2),
											SymbolCount: proto.Int32(3),

											AwardGridVec: []*message.Csh_Grid{
												{
													Row: proto.Int32(2),
												},
												{
													Column: proto.Int32(1),
													Row:    proto.Int32(1),
												},
												{
													Column: proto.Int32(2),
													Row:    proto.Int32(2),
												},
												{
													Column: proto.Int32(2),
													Row:    proto.Int32(3),
												},
											},
										},
									},
								},
								{
									AwardMoney: proto.Float64(1.6),
									AwardSet: []*message.Csh_AwardData{
										{
											AwardMoney:  proto.Float64(1.6),
											AwardSymbol: proto.Int32(3),
											Lines:       proto.Int32(4),
											SymbolCount: proto.Int32(3),
											AwardGridVec: []*message.Csh_Grid{
												{Row: proto.Int32(2)},
												{Row: proto.Int32(3)},
												{Column: proto.Int32(1), Row: proto.Int32(1)},
												{Column: proto.Int32(2)},
												{Column: proto.Int32(2), Row: proto.Int32(2)},
											},
										},
									},
								},
							},
							ColumnSymbol: []*message.Csh_Column{
								{
									Row: []int32{9, 6, 7, 3, 5, 3, 7, 9},
								},
								{
									Row: []int32{0, 8, 9, 9, 3, 3},
								},
								{
									Row: []int32{6, 0, 7, 7, 3, 3, 9, 7, 5},
								},
								{
									Row: []int32{9, 1, 9, 9},
								},
								{
									Row: []int32{9, 2, 9, 5},
								},
								{
									Row: []int32{3, 3, 0, 6},
								},
							},
						},
					},
				},
			},
		}
		encode, _ := proto.Marshal(&spinAll)

		var resData = message.Csh_ResData{
			Type:  proto.Int32(AckType["spin"]),
			Token: proto.String(token),
			Data: []*message.Csh_InfoData{
				{
					Encode: encode,
				},
			},
		}

		*ret, _ = proto.Marshal(&resData)

		return nil
	})
}
