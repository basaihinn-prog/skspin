package fileutil

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"mime/multipart"
	"os"
	"strconv"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

type Land struct {
	LanName      string   `json:"landName"`
	Abbreviation string   `json:"abbreviation"`
	Id           []string `json:"id"`
	Speak        []string `json:"speak"`
}

func CreateAndInputExcel(path string, lan []Land) (error, error) {

	//f := excelize.NewFile()
	//index, _ := f.NewSheet("Sheet1")
	//
	//for i, value := range lan {
	//	f.SetCellValue("Sheet1", Getclom(i+1)+"1", value.LanName)
	//	f.SetCellValue("Sheet1", Getclom(i+1)+"2", value.Abbreviation)
	//}
	//
	//f.SetActiveSheet(index)
	//
	//if err := f.SaveAs(path); err != nil {
	//	return err, errors.New("文件创建失败")
	//}
	return nil, nil
}

func ReadExcel(multipartFile multipart.File) (landList []map[string]string, err error) {
	xlsx, err := excelize.OpenReader(multipartFile)
	if err != nil {
		fmt.Printf("open excel error:[%s]", err.Error())
		return nil, err
	}
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

	for i, row := range rows {
		if i == 0 {
			continue
		}
		cellItem := make(map[string]string)
		for k, v := range row {
			if rows[0][k] != "" {
				cellItem[rows[0][k]] = v
			}

		}
		landList = append(landList, cellItem)
	}

	return
}

func GetExcelCell(path string, sheet string, cell string) (vl string, err error) {

	//f, err := excelize.OpenFile(path)
	//
	//vl, err = f.GetCellValue(sheet, cell)

	return cell, err
}

func Getclom(n int) string {
	if n < 1 {
		return ""
	}
	result := ""
	for n > 0 {
		n--
		remainder := n % 26
		result = strconv.Itoa('A'+remainder) + result
		n /= 26
	}
	return result
}
