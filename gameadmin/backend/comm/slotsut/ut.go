package slotsut

import (
	"crypto/rand"
	"fmt"
	"game/duck/logger"
	"log"
	"math/big"
	"regexp"
	"strings"
)

// rand num is in [0,max)
func UtilRand(max int64) int {
	if max <= 0 {
		panic("rand value not less equal zero")
	}
	n, e := rand.Int(rand.Reader, big.NewInt(max))
	if e != nil {
		log.Fatal(e)
	}
	return int(n.Int64())
}
func Rand(max int64) int {
	return UtilRand(max)
}

func SliceContain(arr []int, d int) bool {
	for _, v := range arr {
		if v == d {
			return true
		}
	}
	return false
}

/*
arr jiu就是权重配置
allWeitht 所有权重sum
return index //返回权重索引
*/
func RandByWeight(arr []int, allWeitht int) int {
	r := UtilRand(int64(allWeitht))
	length := len(arr)
	var index int
	var preV int
	for i := length - 1; i >= 0; i-- {
		allWeitht -= preV
		if r < allWeitht {
			index = i
		} else {
			break
		}
		preV = arr[i]
	}
	return index
}

// [min,max)
func UtilRandDuring(min, max int) int {
	return min + UtilRand(int64(max-min))
}

func ShuffleSlice(strs []string) {
	length := len(strs)
	for i := 0; i < length; i++ {
		r := UtilRand(int64(length))
		strs[i], strs[r] = strs[r], strs[i]
	}
}
func ShuffleIntSlice(arr []int) {
	length := len(arr)
	for i := 0; i < length; i++ {
		r := UtilRand(int64(length))
		arr[i], arr[r] = arr[r], arr[i]
	}
}

func RandPickSome(origins []int, pickCount int) []int {
	newOrigins := make([]int, len(origins), len(origins))
	copy(newOrigins, origins)
	origins = newOrigins
	originLen := len(origins)
	if originLen < pickCount {
		return nil
	} else if originLen == pickCount {
		return origins
	} else {
		results := make([]int, pickCount)
		for i := 0; i < pickCount; i++ {
			ri := UtilRand(int64(len(origins)))
			results[i] = origins[ri]
			origins = append(origins[:ri], origins[ri+1:]...)
		}
		return results
	}
}

func RandPickSomeChangeOrigins(origins *[]int, pickCount int) (result []int) {
	originLen := len(*origins)
	if originLen < pickCount {
		return []int{}
	} else if originLen == pickCount {
		result = *origins
		*origins = []int{}
		return
	} else {
		result = make([]int, pickCount)
		for i := 0; i < pickCount; i++ {
			ri := UtilRand(int64(len(*origins)))
			result[i] = (*origins)[ri]
			*origins = append((*origins)[:ri], (*origins)[ri+1:]...)
		}
		return
	}
}

func ChangeRightkey(c float64) string {
	re, err := regexp.Compile("\\.|-")
	if err != nil {
		logger.Err(err)
	}
	return re.ReplaceAllString(fmt.Sprintf("%v", c), "_")
}
func ChangeDBRightkey(c int) string {
	re, err := regexp.Compile("\\.|-")
	if err != nil {
		logger.Err(err)
	}
	return re.ReplaceAllString(fmt.Sprintf("%v", c), "_")
}

func MinInt(xs ...int) int {

	if len(xs) == 0 {
		logger.Fatal("数组不合法")
	}
	v := xs[0]
	for _, tv := range xs[1:] {
		if v > tv {
			v = tv
		}
	}
	return v
}
func MakeArrAllEqualValue(min, max, length, total int) []int {
	logger.Info(min, max, length, total)
	if min*length > total {
		logger.Fatal("最小间距生成过大")
	}
	if max*length < total {
		logger.Fatal("最大间距生成过小")
	}
	leave := total - min*length
	sp := total / length
	arr := make([]int, length)
	arrIndex := make([]int, length)
	for i := range arr {
		arr[i] = min
		arrIndex[i] = i
	}
	for _, indexV := range arrIndex {
		r := UtilRand(int64(sp + 1 - min))
		leave -= r
		arr[indexV] += r
	}
	ShuffleIntSlice(arrIndex)
	for _, indexV := range arrIndex {
		r := UtilRand(int64(MinInt(max-arr[indexV], leave)) + 1)
		leave -= r
		arr[indexV] += r
	}
	return arr
}

func MakeSlotsDBPath(file string) string {
	if strings.HasPrefix(file, "slotsDB/") {
		return file
	}
	return fmt.Sprintf("slotsDB/%v", file)
}

// func BinaryFind(arr slotsdomain.TimesIndexs, leftIndex int, rightIndex int, findTimes float64, asideIndex *int) {
// 	if leftIndex > rightIndex {
// 		return
// 	}
// 	middle := (leftIndex + rightIndex) / 2
// 	*asideIndex = middle
// 	if (arr)[middle].Times < findTimes {
// 		if rightIndex > middle && arr[middle+1].Times > findTimes {
// 			return
// 		}
// 		BinaryFind(arr, middle+1, rightIndex, findTimes, asideIndex)
// 	} else if (arr)[middle].Times > findTimes {
// 		BinaryFind(arr, leftIndex, middle-1, findTimes, asideIndex)
// 	}
// }
