package PSF_ON_00004

import (
	"serve/fish_comm/rng"
	"serve/service_fish/domain/file"
	"strconv"
)

var FishPath = &fishPath{
	FishInfo: file.FishInfo{
		IsFishPath: true,
		Model:      make(map[string]interface{}),
		FileInfo: file.FileInfo{
			Docker: "/data/script",
			Folder: "domain/fish/PSF-ON-00004",
			File:   "PSF-ON-00004-fishpath.json",
		},
	},
}

type fishPath struct {
	file.FishInfo
}

func init() {
	FishPath.FishInfo.Deserialization()
}

func (f *fishPath) Data(pathId int) []interface{} {
	data := f.Model["fishPath"].(map[string]interface{})

	switch pathId {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		pathData := data["PA0"+strconv.Itoa(pathId)].(map[string]interface{})

		if int(pathData["Direction"].(float64)) == 1 {
			p := pathData["Path"].([]interface{})
			return []interface{}{p[6], p[7], p[4], p[5], p[2], p[3], p[0], p[1]}
		}
		return pathData["Path"].([]interface{})

	case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		fallthrough
	case 20, 21, 22, 23, 24, 25, 26, 27, 28, 29:
		fallthrough
	case 30, 31, 32, 33, 34, 35, 36, 37, 38, 39:
		fallthrough
	case 40, 41, 42, 43, 44, 45, 46, 47, 48, 49:
		fallthrough
	case 50, 51, 52, 53, 54, 55, 56, 57, 58, 59:
		fallthrough
	case 60, 61, 62, 63, 64, 65, 66, 67, 68, 69:
		pathData := data["PA"+strconv.Itoa(pathId)].(map[string]interface{})

		if int(pathData["Direction"].(float64)) == 1 {
			p := pathData["Path"].([]interface{})
			return []interface{}{p[6], p[7], p[4], p[5], p[2], p[3], p[0], p[1]}
		}
		return pathData["Path"].([]interface{})

	case 70:
		return f.Data(f.pickPath(data["RAN_PA00"].(map[string]interface{})))

	case 71:
		return f.Data(f.pickPath(data["RAN_PA01"].(map[string]interface{})))

	case 72:
		return f.Data(f.pickPath(data["RAN_PA02"].(map[string]interface{})))

	case 73:
		return f.Data(f.pickPath(data["RAN_PA03"].(map[string]interface{})))

	case 74:
		return f.Data(f.pickPath(data["RAN_PA04"].(map[string]interface{})))

	case 75:
		return f.Data(f.pickPath(data["RAN_PA05"].(map[string]interface{})))

	case 76:
		return f.Data(f.pickPath(data["RAN_PA06"].(map[string]interface{})))

	case 77:
		return f.Data(f.pickPath(data["RAN_PA07"].(map[string]interface{})))

	case 78:
		return f.Data(f.pickPath(data["RAN_PA08"].(map[string]interface{})))

	default:
		panic("error pathId")
	}

	return nil
}

func (f *fishPath) pickPath(randomData map[string]interface{}) (pathId int) {
	options := make([]rng.Option, 0, len(randomData))

	for i := 1; i < len(randomData); i++ {
		data := randomData["RandomPath_"+strconv.Itoa(i)].(map[string]interface{})
		options = append(options, rng.Option{
			int(data["Weight"].(float64)),
			int(data["PathID"].(float64)),
		})
	}

	return f.Rng(options).(int)
}
