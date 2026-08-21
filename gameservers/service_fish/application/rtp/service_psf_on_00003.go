package rtp

const (
	PSF_ON_00003_MERCENARY_BULLET_LEVEL_1 = 60
	PSF_ON_00003_MERCENARY_BULLET_LEVEL_2 = 105
	PSF_ON_00003_MERCENARY_BULLET_LEVEL_3 = 150
	PSF_ON_00003_MERCENARY_SHOTGUN        = 3
	PSF_ON_00003_MERCENARY_BAZOOKA        = 5

	PSF_ON_00003_MERCENARY_TYPE_RIFLE   = 0
	PSF_ON_00003_MERCENARY_TYPE_SHOTGUN = 1
	PSF_ON_00003_MERCENARY_TYPE_BAZOOKA = 2
)

func setMercenaryDataMap(size int) map[uint64]uint64 {
	data := make(map[uint64]uint64, size)

	for i := 0; i < size; i++ {
		data[uint64(i)] = 0
	}

	return data
}
