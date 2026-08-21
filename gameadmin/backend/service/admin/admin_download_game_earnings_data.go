package main

import (
	"encoding/json"
	"fmt"
	"game/comm"
	"game/duck/lang"
	"io"
	"net/http"

	"github.com/tealeg/xlsx"
)

func DownGameEarningsData(w http.ResponseWriter, r *http.Request) {
	//path := NewPath(r.URL.Path, "admin")
	token := r.Header.Get("Authorization")
	ctx := &Context{
		Request:  r,
		PID:      1,
		Username: "admin",
	}
	cker := &checker{}
	// if strings.HasPrefix(token, "user:") {
	// 	username := strings.TrimPrefix(token, "user:")
	// 	one := adminpb.DBAdminer{}
	// 	err := CollAdminUser.FindOne(bson.M{"Username": username, "Status": comm.User_Status_Normal}, &one)
	// 	if err != nil {
	// 		return
	// 	}
	// 	ctx.PID = one.ID
	// 	ctx.Username = username
	// } else {
	pid, username := cker.GetPidAndUname(token)
	if pid == 0 {
		return
	}
	ctx.PID = pid
	ctx.Username = username
	// }
	params := getGameEarningsDataParams{
		Game:       "",
		OperatorId: 0,
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return
	}
	gameData := &getGameEarningsDataResults{}
	err = getGameEarningsData(ctx, params, gameData)
	if err != nil {
		return
	}
	fmt.Println(gameData)
	fileName := fmt.Sprintf("%s.xlsx", params.Game)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet(params.Game)
	title1 := []string{lang.Get(ctx.Lang, "项目"),
		lang.Get(ctx.Lang, "7日数据"),
		lang.Get(ctx.Lang, "当月数据"),
		lang.Get(ctx.Lang, "总数据")}
	row := sheet.AddRow()
	var cell *xlsx.Cell
	// title
	titleType := xlsx.NewStyle()
	titleType.Fill = *xlsx.NewFill("solid", "999999", "999999")
	titleType.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	for _, title := range title1 {
		cell = row.AddCell()
		cell.SetStyle(titleType)
		cell.Value = title
	}

	gameDataInfos := make([][]string, 3)
	gameDataInfos[0] = []string{
		lang.Get(ctx.Lang, "总产出"),
		fmt.Sprintf("%v", gameData.SevenData.WinAmount),
		fmt.Sprintf("%v", gameData.MonthData.WinAmount),
		fmt.Sprintf("%v", gameData.TotalData.WinAmount),
	}
	gameDataInfos[1] = []string{
		lang.Get(ctx.Lang, "总流水"),
		fmt.Sprintf("%v", gameData.SevenData.BetAmount),
		fmt.Sprintf("%v", gameData.MonthData.BetAmount),
		fmt.Sprintf("%v", gameData.TotalData.BetAmount),
	}
	gameDataInfos[2] = []string{
		lang.Get(ctx.Lang, "总回报率"),
		fmt.Sprintf("%.4f", comm.DIV(float64(gameData.SevenData.WinAmount), float64(gameData.SevenData.BetAmount+gameData.SevenData.WinAmount))),
		fmt.Sprintf("%.4f", comm.DIV(float64(gameData.MonthData.WinAmount), float64(gameData.MonthData.BetAmount+gameData.MonthData.WinAmount))),
		fmt.Sprintf("%.4f", comm.DIV(float64(gameData.TotalData.WinAmount), float64(gameData.TotalData.BetAmount+gameData.TotalData.WinAmount))),
	}

	for i := 0; i < len(gameDataInfos); i++ {
		row = sheet.AddRow()
		for j := 0; j < len(gameDataInfos[i]); j++ {
			row.AddCell().Value = gameDataInfos[i][j]
		}
	}

	for i := 0; i < 5; i++ {
		row = sheet.AddRow()
	}

	row = sheet.AddRow()
	dataT := []string{
		lang.Get(ctx.Lang, "日期"),
		lang.Get(ctx.Lang, "DAU"),
		lang.Get(ctx.Lang, "总压分"),
		lang.Get(ctx.Lang, "系统收益"),
		lang.Get(ctx.Lang, "购买小游戏总押注"),
		lang.Get(ctx.Lang, "购买小游戏回报率"),
		lang.Get(ctx.Lang, "回报率")}
	for i := 0; i < len(dataT); i++ {
		cell = row.AddCell()
		cell.SetStyle(titleType)
		cell.Value = dataT[i]
	}
	gameData.WeekData.NowWeek = append(gameData.WeekData.NowWeek, gameData.WeekData.LastWeek...)
	for i := len(gameData.WeekData.NowWeek) - 1; i >= 0; i-- {
		row = sheet.AddRow()
		row.AddCell().Value = gameData.WeekData.NowWeek[i].Time
		row.AddCell().Value = fmt.Sprintf("%v", gameData.WeekData.NowWeek[i].DauValue)
		row.AddCell().Value = fmt.Sprintf("%v", gameData.WeekData.NowWeek[i].BetValue)
		row.AddCell().Value = fmt.Sprintf("%v", gameData.WeekData.NowWeek[i].SystemValue)
		row.AddCell().Value = fmt.Sprintf("%v", gameData.WeekData.NowWeek[i].BuyBet)
		row.AddCell().Value = fmt.Sprintf("%v", gameData.WeekData.NowWeek[i].BuyRateValue)
		row.AddCell().Value = fmt.Sprintf("%.4f", gameData.WeekData.NowWeek[i].RateValue)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	file.Write(w)
}
