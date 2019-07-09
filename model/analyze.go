package model

import (
	"fmt"
	"gw/utils"
	"time"
)

type ReqAnalyze struct {
	StartTime int `json:"start_time"`
	CompanyId int `json:"companyId"`
}

func (b *Model) GetAlertAnalyze(params ReqAnalyze, companyId int) ([]GwAnalyze, error) {
	var (
		analyze        []GwAnalyze
		queryCompanyId int
	)

	db, _ := utils.Connect()
	db.ShowSQL(true)

	queryCompanyId = params.CompanyId
	if companyId != 1 {
		queryCompanyId = companyId
	}

	if params.StartTime == 0 {
		var (
			timeMap []map[string]string
		)
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL 3 MONTH), '%Y-%m-01 00:00:00') AS 'startTime'")
		timeTemplate := "2006-01-02 15:04:05" //常规类型
		timeStart, _ := time.ParseInLocation(timeTemplate, timeMap[0]["startTime"], time.Local)
		fmt.Println(timeMap)
		params.StartTime = int(timeStart.Unix())
	}
	fmt.Println(params.StartTime)
	if err := db.Where("daytime >= ?", params.StartTime).And("company_id = ?", queryCompanyId).Find(&analyze); err != nil {
		return analyze, err
	}
	return analyze, nil
}
