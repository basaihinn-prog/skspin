//go:build dev || staging
// +build dev staging

package PSF_ON_00007_1_MONSTER

func (l *luckSlot) rtpHit(
	rtp *struct{ HitBsLuckSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink *FsChangeLink1,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	pirateShip, treasureDish, giantWhale *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := l.hit(fsMath, fsStripMath, changeLink,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			pirateShip, treasureDish, giantWhale,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
