package controllers

import (
	"strconv"
	"strings"
	"time"

	"fmt"
	"gw/model"

	"gw/config"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type reqAlertInfo struct {
	DeviceId    string `json:"device_id" binding: "required"`
	TokenString string `json:"token" binding:"required"`
}

type reqDeviceHistory struct {
	TokenString string `json:"token" binding:"required"`
	SelectTime  int    `json:"selectTime" binding:"required"`
	DeviceId    string `json:"device_id" binding:"required"`
}

type reqDeviceOPenRecor struct {
	Token
	DeviceId string `json:"device_id" binding:"required"`
}
type reqUserNotifyHistory struct {
	Token
	UserId   int    `json:"user_id" binding:"required"`
	DeviceId string `json:"device_id" binding:"required"`
}

func (tc *BaseController) AlertList(c *gin.Context) {
	var (
		token model.DeviceParams
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, comId, autherr := tc.CheckAuth(token.TokenString, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	count, alert, err := m.GetAlertList(token, comId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data := struct {
		Count int64                 `json:"count"`
		List  []model.GetDeviceList `json:"list"`
	}{
		count,
		alert,
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) AlertInfo(c *gin.Context) {
	var (
		token reqAlertInfo
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, "", true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	info, err := m.GetDeviceAlertInfoByDeviceId(token.DeviceId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": info})
	return
}

func (tc *BaseController) DevicePressureHistory(c *gin.Context) {
	var (
		token reqDeviceHistory
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, "", true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	data, err := m.GetDevicePressureHistory(token.DeviceId, token.SelectTime)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceOpenHistory(c *gin.Context) {
	var (
		token reqDeviceOPenRecor
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, "", true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data, err := m.GetDeviceOpenRecord(token.DeviceId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) UserNotifyHistory(c *gin.Context) {
	var (
		token reqUserNotifyHistory
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, "", true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data, err := m.GetUserNotifyHistory(token.UserId, token.DeviceId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) AlertHistory(c *gin.Context) {
	var (
		token model.ReqAlertHistory
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, comId, autherr := tc.CheckAuth(token.TokenString, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	total, list, err := m.GetAlertHistory(token, comId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data := struct {
		List  []model.ResAlertForHistory `json:"list"`
		Total int64                      `json:"total"`
	}{
		list,
		total,
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceAlertEvent(c *gin.Context) {
	var (
		token model.ReqAlertEvent
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, _, autherr := tc.CheckAuth(token.TokenString, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	num, list, err := m.GetDeviceAlertEvent(token)
	data := struct {
		List  []model.ResDeviceAlertEvent `json:"list"`
		Total int64                       `json:"total"`
	}{
		list,
		num}
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceAlertOriginData(c *gin.Context) {
	var (
		token model.ReqAlertEventOriginData
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, _, autherr := tc.CheckAuth(token.TokenString, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	num, list, err := m.GetDeviceAlertOriginData(token)
	data := struct {
		List  []model.GwAlert `json:"list"`
		Total int64           `json:"total"`
	}{
		list,
		num}
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceAlertDetail(c *gin.Context) {
	var (
		token model.ReqDeviceAlertDetail
	)

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	if _, _, err := tc.CheckAuth(token.TokenString, "", true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	total, list, err := m.GetDeviceAlertDetail(token)
	data := struct {
		List  []model.ResAlertDetail `json:"data"`
		Total int64                  `json:"total"`
	}{
		list,
		total}

	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceEventExport(c *gin.Context) {
	var (
		showType   string
		timeType   int
		dataType   int
		selectTime int
	)
	router := c.Request.RequestURI
	token := c.PostForm("token")
	exportType := c.PostForm("exportType")
	deviceId := c.PostForm("deviceId")

	if _, _, err := tc.CheckAuth(token, router, false); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	timeLayout := "2006-01-02 15:04:05" //转化所需模板

	if exportType == "0" {
		showType = c.PostForm("showType")
		eventType, _ := strconv.Atoi(showType)
		list, err := m.ExportDeviceAlertEvent(eventType, deviceId)
		if err != nil {
			c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
			return
		}
		fmt.Println(list)
		// f := excelize.NewFile()
		// index := f.NewSheet("Sheet1")
		// timeLayout := "2006-01-02 15:04:05" //转化所需模板
		for key, val := range list {
			i := key + 1
			f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i), val.Id)
			senttime64, _ := strconv.ParseInt(val.Sendtime, 10, 64)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), time.Unix(senttime64, 0).Format(timeLayout))
			if val.AlertType == "20" {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "阀门打开")
			} else if val.AlertType == "30" {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "撞倒")
			} else if val.AlertType == "60" {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "水压异常")
			} else if val.AlertType == "70" {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "失联")
			}
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i), val.Cola)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i), val.CompanyName)
			f.SetCellValue("Sheet1", fmt.Sprintf("F%v", i), val.Descrip)
		}

	} else if exportType == "1" {
		timeType, _ = strconv.Atoi(c.PostForm("timeType"))
		dataType, _ = strconv.Atoi(c.PostForm("dataType"))
		selectTime, _ = strconv.Atoi(c.PostForm("selectTime"))
		list, err := m.ExportDeviceAlertOriginData(timeType, selectTime, dataType, deviceId)
		if err != nil {
			c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
			return
		}

		for key, val := range list {
			i := key + 1
			senttime64, _ := strconv.ParseInt(val.Sendtime, 10, 64)
			f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i), time.Unix(senttime64, 0).Format(timeLayout))
			if val.AlertType == "20" {
				f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), "阀门打开")
			} else if val.AlertType == "30" {
				f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), "撞倒")
			} else if val.AlertType == "60" {
				f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), "水压异常")
			} else if val.AlertType == "70" {
				f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), "失联")
			}
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), val.Cola)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i), "已解析")
		}
	}

	f.SetActiveSheet(index)
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".xlsx"
	fileerr := f.SaveAs(config.ExportFolder + fileName)
	if fileerr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": fileerr.Error()})
		return
	}
	extraHeaders := make(map[string]string)
	extraHeaders["Content-Disposition"] = `attachment; filename="` + fileName + `"`
	extraHeaders["Expires"] = "0"
	extraHeaders["Cache-Control"] = "must-revalidate"
	extraHeaders["Pragma"] = "public"
	nf, _ := os.Open(config.ExportFolder + fileName)
	defer nf.Close()
	fi, _ := nf.Stat()
	c.DataFromReader(200, fi.Size(), "application/octet-stream", nf, extraHeaders)
	return
}

func (tc *BaseController) ExportAlertTrace(c *gin.Context) {
	token := c.PostForm("token")
	selectCompanyId, _ := strconv.Atoi(c.PostForm("companyId"))
	addkeys := c.PostForm("addkeys")
	alertStae, _ := strconv.Atoi(c.PostForm("alertState"))
	dataString := c.PostForm("dataPicker")
	dataPicker := strings.Split(dataString, ",")
	_, companyId, autherr := tc.CheckAuth(token, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	if companyId != 1 && companyId != selectCompanyId {
		c.JSON(200, gin.H{"status": -1, "msg": fmt.Errorf("You have no permisson!")})
		return
	}
	list, err := m.ExportAlertTrace(selectCompanyId, addkeys, alertStae, dataPicker)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	timeLayout := "2006-01-02 15:04:05" //转化所需模板

	for key, val := range list {
		i := key + 1
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i), val.CompanyName)
		senttime64, _ := strconv.ParseInt(val.Sendtime, 10, 64)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), time.Unix(senttime64, 0).Format(timeLayout))
		if val.AlertType == "20" {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "阀门打开")
		} else if val.AlertType == "30" {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "撞倒")
		} else if val.AlertType == "60" {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "水压异常")
		} else if val.AlertType == "70" {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "失联")
		}
		if val.Restoretime != 0 {
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i), "已解除")
		} else {
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i), "未解除")
		}
		if val.NotifyStatus == 1 {
			f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i), "通知已到达")
		} else {
			f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i), "通知未到达")
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("F%v", i), val.DeviceId)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%v", i), val.Address)
	}

	f.SetActiveSheet(index)
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".xlsx"
	fileerr := f.SaveAs(config.ExportFolder + fileName)
	if fileerr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": fileerr.Error()})
		return
	}
	extraHeaders := make(map[string]string)
	extraHeaders["Content-Disposition"] = `attachment; filename="` + fileName + `"`
	extraHeaders["Expires"] = "0"
	extraHeaders["Cache-Control"] = "must-revalidate"
	extraHeaders["Pragma"] = "public"
	nf, _ := os.Open(config.ExportFolder + fileName)
	defer nf.Close()
	fi, _ := nf.Stat()
	c.DataFromReader(200, fi.Size(), "application/octet-stream", nf, extraHeaders)
	return
}
