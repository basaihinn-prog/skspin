package PSF_ON_00001

import (
	"log"
	"serve/service_fish/domain/file"
	"strconv"
)

var Objects = &objects{
	FishInfo: file.FishInfo{
		IsFishPath: false,
		Model:      make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00001",
			File:   "PSF-ON-00001-Objects.json",
		},
	},
}

type objects struct {
	file.FishInfo
}

func init() {
	Objects.FishInfo.Deserialization()
}

func (o *objects) data(objectId int) map[string]interface{} {
	data := o.Model["Objects"].(map[string]interface{})

	switch objectId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return data["F0"+strconv.Itoa(objectId)].(map[string]interface{})

	case 30:
		return data["F30"].(map[string]interface{})

	case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 26, 27, 28, 29:
		return data["F"+strconv.Itoa(objectId)].(map[string]interface{})

	default:
		log.Println("Objects error groupId", objectId)
		panic("Objects error groupId")
	}
}

func (o *objects) Speed(objectId int) (fishSpeed int32) {
	return int32(o.data(objectId)["Speed"].(float64))
}

func (o *objects) Layer(objectId int) (fishLayer int32) {
	return int32(o.data(objectId)["Layer"].(float64))
}

func (o *objects) ExtraData(objectId int) (extraData []interface{}) {
	return o.data(objectId)["ExtraData"].([]interface{})
}
