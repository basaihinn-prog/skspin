package comm

import "time"

// 获取周一的零点
func GetMondayOfWeek() time.Time {
	t := time.Now()
	dayObj := GetZeroTime(t)
	if t.Weekday() == time.Monday {
		return dayObj
	}
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset -= 6
	}
	return dayObj.AddDate(0, 0, offset)
}

// 获取上一周周一零点
func GetLastMondayOfWeek() time.Time {
	return GetMondayOfWeek().AddDate(0, 0, -7)
}

// 获取某一天零点
func GetZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
