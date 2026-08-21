package probability

import "github.com/google/uuid"

type Probability struct {
	FishTypeId        int
	Pay               int
	Bullet            int
	BonusPayload      interface{}
	TriggerIconId     int
	BonusTypeId       int
	Multiplier        int
	BonusUuid         string
	ExtraData         []interface{}
	ExtraTriggerBonus []*Probability
}

func newSimplePay(fishTypeId, pay, multiplier int) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: -1,
		BonusTypeId:   -1,
		Multiplier:    multiplier,
		BonusPayload:  -1,
		BonusUuid:     "",
	}
}

func newSimpleBullet(fishTypeId, pay, bullet, triggerIconId, bonusTypeId, multiplier int) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		Bullet:        bullet,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		BonusPayload:  -1,
		BonusUuid:     uuid.New().String(),
	}
}

func newTriggerBonusPay(fishTypeId, pay, triggerIconId, bonusTypeId, multiplier int) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		BonusPayload:  -1,
		BonusUuid:     uuid.New().String(),
	}
}

func newOptionBonusPay(fishTypeId, pay, triggerIconId, bonusTypeId, multiplier int, bonusPayload interface{}) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		BonusPayload:  bonusPayload,
		BonusUuid:     uuid.New().String(),
	}
}

func newOptionBonusBullet(fishTypeId, pay, triggerIconId, bonusTypeId, multiplier, bullet int, bonusPayload interface{}) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		Bullet:        bullet,
		BonusPayload:  bonusPayload,
		BonusUuid:     uuid.New().String(),
	}
}

func newExtraOptionBonusPay(fishTypeId, pay, triggerIconId, bonusTypeId, multiplier int, extraData []interface{}, bonusPayload interface{}) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		BonusPayload:  bonusPayload,
		BonusUuid:     uuid.New().String(),
		ExtraData:     extraData,
	}
}

func newExtraOptionBonusBullet(fishTypeId, pay, triggerIconId, bonusTypeId, multiplier, bullets int, extraData []interface{}, bonusPayload interface{}) *Probability {
	return &Probability{
		FishTypeId:    fishTypeId,
		Pay:           pay,
		TriggerIconId: triggerIconId,
		BonusTypeId:   bonusTypeId,
		Multiplier:    multiplier,
		Bullet:        bullets,
		BonusPayload:  bonusPayload,
		BonusUuid:     uuid.New().String(),
		ExtraData:     extraData,
	}
}
