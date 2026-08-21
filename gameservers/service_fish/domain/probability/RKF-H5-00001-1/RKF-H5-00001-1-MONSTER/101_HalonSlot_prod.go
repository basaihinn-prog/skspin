//go:build prod
// +build prod

package RKF_H5_00001_1_MONSTER

func (h *halonSlot) rtpHit(
	rtp *struct{ HitBsHalonSlotMath },
	fsMath *FsSlotMath,
	fsStripMath *FsStripTableMath,
	changeLink1Math *FsChangeLink1Math, changeLink2Math *FsChangeLink2Math,
	fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
	faceOff1, faceOff3, faceOff5 *FsFishMath,
) (iconPay, triggerIconId, bonusTypeId int, iconIds []int32, iconPays []int64) {
	triggerIconId = -1
	bonusTypeId = -1

	if MONSTER.isHit(rtp.TriggerWeight) {
		triggerIconId = rtp.TriggerIconID
		bonusTypeId = rtp.Type
	}

	if MONSTER.isHit(rtp.HitWeight) {
		v1, v2, v3 := h.hit(
			fsMath, fsStripMath,
			changeLink1Math, changeLink2Math,
			fish0, fish4, fish9, fish11, fish13, fish14, fish15, fish16, fish17,
			faceOff1, faceOff3, faceOff5,
		)
		return v1, triggerIconId, bonusTypeId, v2, v3
	}

	return 0, triggerIconId, bonusTypeId, nil, nil
}
