package PSF_ON_00005

import (
	"log"
	"serve/fish_comm/rng"
	"serve/service_fish/domain/file"
	"strconv"
)

var Groups = &groups{
	FishInfo: file.FishInfo{
		IsFishPath: false,
		Model:      make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00005",
			File:   "PSF-ON-00005-Groups.json",
		},
	},
}

type groups struct {
	file.FishInfo
}

func init() {
	Groups.FishInfo.Deserialization()
}

func (g *groups) Data(groupId int) (fishGroup map[string]interface{}, fishId int) {
	data := g.Model["Groups"].(map[string]interface{})

	switch groupId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		if v, ok := data["GR0"+strconv.Itoa(groupId)]; ok {
			return v.(map[string]interface{}), -1
		}
		return nil, -1

	case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		if v, ok := data["GR"+strconv.Itoa(groupId)]; ok {
			return v.(map[string]interface{}), -1
		}
		return nil, -1

	case 24:
		if v, ok := data["RAN_F00"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
		}
		return nil, -1

	case 25:
		if v, ok := data["RAN_F01"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
		}
		return nil, -1

	case 26:
		if v, ok := data["RAN_F02"]; ok {
			randomData := v.(map[string]interface{})
			return nil, g.pickFish(randomData["RandomFish"].(map[string]interface{}))
		}
		return nil, -1

	case 27:
		if v, ok := data["RAN_F03"]; ok {
			randomData := v.(map[string]interface{})
			return nil, g.pickFish(randomData["RandomFish"].(map[string]interface{}))
		}
		return nil, -1

	default:
		log.Println("Groups error groupId", groupId)
		return nil, -1
	}
}

func (g *groups) pickGroup(randomData map[string]interface{}) (groupId int) {
	options := make([]rng.Option, 0, len(randomData))

	for i := 1; i < len(randomData); i++ {
		data := randomData["RandomFish_"+strconv.Itoa(i)].(map[string]interface{})
		options = append(options, rng.Option{
			int(data["Weight"].(float64)),
			int(data["GroupID"].(float64)),
		})
	}

	return g.Rng(options).(int)
}

func (g *groups) pickFish(randomData map[string]interface{}) (FishId int) {
	options := make([]rng.Option, 0, len(randomData))

	for i := 1; i < len(randomData); i++ {
		data := randomData["RandomFish_"+strconv.Itoa(i)].(map[string]interface{})
		options = append(options, rng.Option{
			int(data["Weight"].(float64)),
			int(data["FishID"].(float64)),
		})
	}

	return g.Rng(options).(int)
}
