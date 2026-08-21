package gamedata

// 函数：判断一个整数是否在切片中
func Contains(slice []int, number int) bool {
	for _, v := range slice {
		if v == number {
			return true
		}
	}
	return false
}
