package PSF_ON_00006

import (
	"log"
	"serve/service_fish/domain/file"
	"strconv"

	"serve/fish_comm/rng"
)

var Groups = &groups{
	FishInfo: file.FishInfo{
		IsFishPath: false,
		Model:      make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00006",
			File:   "PSF-ON-00006-Groups.json",
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

	case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		fallthrough
	case 20, 21, 22, 23, 24, 25, 26, 27, 28, 29:
		fallthrough
	case 30, 31, 32, 33, 34, 35, 36, 37, 38:
		if v, ok := data["GR"+strconv.Itoa(groupId)]; ok {
			return v.(map[string]interface{}), -1
		}
		return nil, -1

	case 39:
		if v, ok := data["RAN_F00"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
		}
		return nil, -1

	case 40:
		if v, ok := data["RAN_F01"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
		}
		return nil, -1

	case 41:
		if v, ok := data["RAN_F02"]; ok {
			randomData := v.(map[string]interface{})
			return nil, g.pickFish(randomData["RandomFish"].(map[string]interface{}))
		}
		return nil, -1

	case 42:
		if v, ok := data["RAN_F03"]; ok {
			randomData := v.(map[string]interface{})
			return nil, g.pickFish(randomData["RandomFish"].(map[string]interface{}))
		}
		return nil, -1

	case 43:
		if v, ok := data["RAN_F04"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
		}
		return nil, -1

	case 44:
		if v, ok := data["RAN_F05"]; ok {
			randomData := v.(map[string]interface{})
			return g.Data(g.pickGroup(randomData["RandomFish"].(map[string]interface{})))
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
func (g *groups) GroupType(groupId int) (groupType int32) {
	groupData, _ := g.Data(groupId)
	return int32(groupData["GroupType"].(float64))
}
