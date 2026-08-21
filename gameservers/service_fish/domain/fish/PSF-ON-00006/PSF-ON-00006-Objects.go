package PSF_ON_00006

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
			Folder: "domain/fish/PSF-ON-00006",
			File:   "PSF-ON-00006-Objects.json",
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

	case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		fallthrough

	case 100, 101, 201, 202, 300, 301, 302, 400:
		return data["F"+strconv.Itoa(objectId)].(map[string]interface{})

	default:
		log.Print("Objects error groupId", objectId)
		panic("Objects error groupId")
	}

	return nil
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
