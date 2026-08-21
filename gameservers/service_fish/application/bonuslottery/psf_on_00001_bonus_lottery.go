package bonuslottery

import (
	"serve/service_fish/application/accountingmanager"
	"serve/service_fish/application/activity"
	"serve/service_fish/application/gamesetting"
	"serve/service_fish/domain/probability"
	probability_proto "serve/service_fish/domain/probability/proto"
	"serve/service_fish/domain/redenvelope"
	"serve/service_fish/domain/slot"
	"time"

	"serve/fish_comm/broadcaster"
	common_proto "serve/fish_comm/common/proto"
	errorcode "serve/fish_comm/error-code"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/wallet"

	"github.com/gogo/protobuf/proto"
)

func psf_on_00001_RedEnvelope(action *flux.Action, depositMultiple uint64) {
	switch action.Key().Name() {
	case redenvelope.ActionRedEnvelopePick:
		re := action.Payload()[1].(*redenvelope.RedEnvelope)

		switch wallet.Service.WalletCategory(re.SecWebSocketKey) {
		case wallet.CategoryHost:
			accountingSn := accountingmanager.Service.GetAccountingSn(re.SecWebSocketKey)
			if accountingSn == 0 {
				balance, _ := wallet.Service.BalanceCache(re.SecWebSocketKey)

				depositCent := re.Bet * re.Rate * depositMultiple
				if balance < depositCent {
					depositCent = balance
				}

				wallet.Service.Deposit(re.SecWebSocketKey, depositCent, accountingSn)
			}

			wallet.Service.IncreaseCent(re.SecWebSocketKey, re.Uuid, re.Bet*re.Pay*re.Rate)
		default:
			wallet.Service.IncreaseBalance(re.SecWebSocketKey, re.Uuid, re.Bet*re.Pay*re.Rate, 0, 0)
		}

		activity.Service.RecordBonus(re.SecWebSocketKey, re)

		// bot will get nil
		if am := accountingmanager.Service.Get(re.SecWebSocketKey); am != nil {
			// TODO JOHNNY Increase 因為 DB IO 會比較慢  先暫時 sleep 暫解 非同步問題
			time.Sleep(500 * time.Millisecond)
			am.PickBonus(re.SecWebSocketKey, re)
		}

		data, dataBuffer := psf_on_00001_RedEnvelope_eOptionRecall(re.SecWebSocketKey, re)

		flux.Send(probability.EMSGID_eOptionRecall, Service.Id, re.SecWebSocketKey, dataBuffer)
		flux.Send(probability.EMSGID_eBroadcastOption, Service.Id, re.RoomUuid, psf_on_00001_eBroadcastOption(data))

	case redenvelope.ActionRedEnvelopeNotExist:
		secWebSocketKey := action.Payload()[0].(string)
		re := action.Payload()[1].(*redenvelope.RedEnvelope)

		logger.Service.Zap.Errorw(BonusLottery_REDENVELOPE_NOT_FOUND,
			"GameUser", re.SecWebSocketKey,
			"GameRoomUuid", re.RoomUuid,
			"FishType", re.FishTypeId,
			"BonusUuid", re.Uuid,
			"BonusPayload", re.BonusPayload,
			"Bet", re.Bet,
			"Rate", re.Rate,
			"Pay", re.Pay,
			"PlayerOptIndex", re.PlayerOptIndex,
		)
		errorcode.Service.Fatal(secWebSocketKey, BonusLottery_REDENVELOPE_NOT_FOUND)
	}
}

func psf_on_00001_RedEnvelope_eOptionRecall(secWebSocketKey string, re *redenvelope.RedEnvelope) (*probability_proto.OptionRecall, []byte) {
	balance, err := wallet.Service.Balance(secWebSocketKey)

	if wallet.Service.WalletCategory(secWebSocketKey) == wallet.CategoryHost {
		cent, err1 := wallet.Service.Cent(secWebSocketKey)

		if err1 != nil {
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"Error", err1,
			)
			panic(secWebSocketKey + err1.Error())
		}

		balance += cent
	} else {
		if err != nil {
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"Error", err,
			)
			panic(secWebSocketKey + err.Error())
		}
	}

	data := &probability_proto.OptionRecall{
		Msgid:          common_proto.EMSGID_eOptionRecall,
		StatusCode:     common_proto.Status_kSuccess,
		RoomUuid:       re.RoomUuid,
		SeatId:         re.SeatId,
		BonusUuid:      re.Uuid,
		PlayerOptIndex: re.PlayerOptIndex,
		HitResult: &probability_proto.HitResult{
			Pay:        re.Bet * re.Pay * re.Rate,
			Multiplier: 1, // not used yet
			PlayerCent: balance,
		},
	}

	for _, v := range re.AllPay {
		data.AllPay = append(data.AllPay, int64(re.Bet)*v*int64(re.Rate))
	}

	flux.Send(broadcaster.ActionBroadcasterBigWin, Service.Id, broadcaster.Service.Id, broadcaster.New(
		secWebSocketKey,
		gamesetting.Service.GameId(secWebSocketKey),
		wallet.Service.MemberId(secWebSocketKey),
		gamesetting.Service.BetList(secWebSocketKey),
		gamesetting.Service.RateList(secWebSocketKey),
		gamesetting.Service.MathModuleId(secWebSocketKey),
		int(gamesetting.Service.RateIndex(secWebSocketKey)),
		re.FishTypeId,
		re.Pay,
		data.HitResult.Pay,
		gamesetting.Service.Rate(secWebSocketKey),
	))

	dataBuffer, _ := proto.Marshal(data)
	return data, dataBuffer
}

func psf_on_00001_Slot(action *flux.Action, depositMultiple uint64) {
	switch action.Key().Name() {
	case slot.ActionSlotPick:
		sl := action.Payload()[1].(*slot.Slot)

		switch wallet.Service.WalletCategory(sl.SecWebSocketKey) {
		case wallet.CategoryHost:
			accountingSn := accountingmanager.Service.GetAccountingSn(sl.SecWebSocketKey)
			if accountingSn == 0 {
				balance, _ := wallet.Service.BalanceCache(sl.SecWebSocketKey)

				depositCent := sl.Bet * sl.Rate * depositMultiple
				if balance < depositCent {
					depositCent = balance
				}

				wallet.Service.Deposit(sl.SecWebSocketKey, depositCent, accountingSn)
			}

			wallet.Service.IncreaseCent(sl.SecWebSocketKey, sl.Uuid, sl.Bet*sl.Pay*sl.Rate)
		default:
			wallet.Service.IncreaseBalance(sl.SecWebSocketKey, sl.Uuid, sl.Bet*sl.Pay*sl.Rate, 0, 0)
		}

		activity.Service.RecordBonus(sl.SecWebSocketKey, sl)

		// bot will get nil
		if am := accountingmanager.Service.Get(sl.SecWebSocketKey); am != nil {
			// TODO JOHNNY Increase 因為 DB IO 會比較慢  先暫時 sleep 暫解 非同步問題
			time.Sleep(500 * time.Millisecond)
			am.PickBonus(sl.SecWebSocketKey, sl)
		}

		data, dataBuffer := psf_on_00001_Slot_eOptionRecall(sl.SecWebSocketKey, sl)

		flux.Send(probability.EMSGID_eOptionRecall, Service.Id, sl.SecWebSocketKey, dataBuffer)
		flux.Send(probability.EMSGID_eBroadcastOption, Service.Id, sl.RoomUuid, psf_on_00001_eBroadcastOption(data))

	case slot.ActionSlotNotExist:
		secWebSocketKey := action.Payload()[0].(string)
		sl := action.Payload()[1].(*slot.Slot)

		logger.Service.Zap.Errorw(BonusLottery_SLOT_NOT_FOUND,
			"GameUser", sl.SecWebSocketKey,
			"GameRoomUuid", sl.RoomUuid,
			"FishType", sl.FishTypeId,
			"BonusUuid", sl.Uuid,
			"BonusPayload", sl.AllPay,
			"Bet", sl.Bet,
			"Rate", sl.Rate,
			"Pay", sl.Pay,
		)
		errorcode.Service.Fatal(secWebSocketKey, BonusLottery_SLOT_NOT_FOUND)
	}
}

func psf_on_00001_Slot_eOptionRecall(secWebSocketKey string, slot *slot.Slot) (*probability_proto.OptionRecall, []byte) {
	balance, err := wallet.Service.Balance(secWebSocketKey)

	if wallet.Service.WalletCategory(secWebSocketKey) == wallet.CategoryHost {
		cent, err1 := wallet.Service.Cent(secWebSocketKey)

		if err1 != nil {
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"Error", err1,
			)
			panic(secWebSocketKey + err1.Error())
		}

		balance += cent
	} else {
		if err != nil {
			logger.Service.Zap.Errorw("BUG",
				"GameUser", secWebSocketKey,
				"Error", err,
			)
			panic(secWebSocketKey + err.Error())
		}
	}

	data := &probability_proto.OptionRecall{
		Msgid:          common_proto.EMSGID_eOptionRecall,
		StatusCode:     common_proto.Status_kSuccess,
		RoomUuid:       slot.RoomUuid,
		SeatId:         slot.SeatId,
		BonusUuid:      slot.Uuid,
		PlayerOptIndex: -1,
		SymbolId:       slot.Reels,
		AllPay:         slot.AllPay,
		HitResult: &probability_proto.HitResult{
			Pay:        slot.Bet * slot.Pay * slot.Rate,
			Multiplier: 1, // not used yet
			PlayerCent: balance,
		},
	}

	flux.Send(broadcaster.ActionBroadcasterBigWin, Service.Id, broadcaster.Service.Id, broadcaster.New(
		secWebSocketKey,
		gamesetting.Service.GameId(secWebSocketKey),
		wallet.Service.MemberId(secWebSocketKey),
		gamesetting.Service.BetList(secWebSocketKey),
		gamesetting.Service.RateList(secWebSocketKey),
		gamesetting.Service.MathModuleId(secWebSocketKey),
		int(gamesetting.Service.RateIndex(secWebSocketKey)),
		slot.FishTypeId,
		slot.Pay,
		data.HitResult.Pay,
		gamesetting.Service.Rate(secWebSocketKey),
	))

	dataBuffer, _ := proto.Marshal(data)
	return data, dataBuffer
}

func psf_on_00001_eBroadcastOption(optionRecall *probability_proto.OptionRecall) []byte {
	data := &probability_proto.BroadcastOption{
		Msgid:        common_proto.EMSGID_eBroadcastOption,
		OptionRecall: optionRecall,
	}
	dataBuffer, _ := proto.Marshal(data)
	return dataBuffer
}
