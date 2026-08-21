package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"serve/servicejili/jili_110_ali/internal/models"

	"serve/comm/plats/platcomm"

	"serve/servicejili/jiliut"

	"github.com/samber/lo"
)

func (f *Fetcher) invokeHistory(path, query string, ret any) (err error) {
	historyHost := "history-api" + f.launchUrl.Host[strings.IndexByte(f.launchUrl.Host, '.'):]
	ul := url.URL{
		Scheme:   "https",
		Host:     historyHost,
		Path:     path,
		RawQuery: query,
	}
	historyRecodeHeaders := map[string]string{
		"Authorization": "Bearer " + f.historyToken,
	}

	var data struct {
		Code    int
		Message string
		Data    any
	}
	data.Data = ret

	lg := slog.With("path", path, "query", query, "user", f.username)
	defer func() {
		lg.With("err", err, "data", data).Info("invokeHistory")
	}()

	err = jiliut.GetJsonWithHeaders(ul.String(), &data, historyRecodeHeaders)
	if err != nil {
		return
	}

	if data.Code != 0 {
		err = errors.New(data.Message)
		return
	}

	return
}

func (f *Fetcher) fetchHisotry(roundIndex string, doc *models.RawSpin) {
	//roundIndex = strings.ReplaceAll(roundIndex, "-", "")

	historyRecodeQuery := lo.Must(url.ParseQuery("EndRowIndex=1&LangId=en-US&LogIndexAsRoundIndex=false&RoundIndex=1718339075110445002&StartRowIndex=1"))
	historyRecodeQuery.Set("RoundIndex", roundIndex)

	var records []*models.HistoryRecord

	time.Sleep(3 * time.Second)
	platcomm.Retry(20, 3*time.Second, func() bool {
		lo.Must0(f.invokeHistory("/history/ali/get-history-record", historyRecodeQuery.Encode(), &records))

		n := len(records)
		lo.Must0(n < 2)
		return n == 0
	})

	doc.HistoryRecord = records[0]

	// https://history-api.jlfafafa3.com/history/csh/get-single-round-log-summary/en-US/1718339075129795002

	var summaries []*models.SingleRoundLogSummary
	lo.Must0(f.invokeHistory("/history/ali/get-single-round-log-summary/en-US/"+roundIndex, "", &summaries))

	doc.SingleRoundLogSummaries = summaries
	doc.LogPlateInfos = make(map[string][]*models.LogPlateInfo, len(summaries))

	for _, summary := range summaries {
		// fmt.Println(summary)

		// https://history-api.jlfafafa3.com/history/csh/get-log-plate-info/1718339075129795002/1718339075129795002

		path := fmt.Sprintf("/history/ali/get-log-plate-info/%s/%s", roundIndex, summary.LogIndex)

		var platInfos []*models.LogPlateInfo
		platcomm.Retry(20, 3*time.Second, func() bool {
			lo.Must0(f.invokeHistory(path, "", &platInfos))
			return len(platInfos) == 0
		})

		lo.Must0(f.invokeHistory(path, "", &platInfos))
		lo.Must0(len(platInfos) != 0)
		doc.LogPlateInfos[summary.LogIndex] = platInfos
	}

}
