package model

import (
	"fmt"
	"gw/utils"

	"github.com/go-xorm/xorm"
)

type DeviceParams struct {
	TokenString string `json:"token" binding:"required"`
	BaseQueryParam
	State string `json:"alarmType"`
	ComId int    `json:"companyId"`
}

type ReqAlertHistory struct {
	TokenString string `json:"token" binding:"required"`
	BaseQueryParam
	DataPicker []string `json:"dataPicker"`
	Addkeys    string   `json:"addkeys"`
	AlertState int      `json:"alertState"`
	CompanyId  int      `json:"companyId"`
}

type ResAlertForHistory struct {
	GwAlert `xorm:"extends"`
	// AlertState   int    `json:"alertState"`
	NotifyStatus int    `json:"notifyStatus"`
	CompanyName  string `json:"companyName"`
	Address      string `json:"address"`
}

type ResAlertHistory struct {
	List  []ResAlertForHistory `json:"data"`
	Total int64                `json:"total"`
}

type DeviceAlertInfo struct {
	Address  string      `json:"address"`
	Cola     string      `json:"cola"`
	State    string      `json:"alarm_type"`
	Value1   string      `json:"value1"`
	Value2   string      `json:"value2"`
	SendTime string      `json:"send_time"`
	DeviceId string      `json:"deviceId"`
	Teles    []GwUser    `json:"teles"`
	Notify   []ResNotify `json:"notify_infos"`
}

type ResNotify struct {
	Data  []GwNotify `json:"data"`
	Phone string     `json:"phone"`
	Name  string     `json:"name"`
}

type DevicePressureHistory struct {
	List   []GwPressure `json:"list"`
	Value1 string       `json:"value1"`
	Value2 string       `json:"value2"`
}

type ResDeviceAlertEvent struct {
	GwAlert     `xorm:"extends"`
	CompanyName string `json:"company"`
}

type ReqAlertEvent struct {
	DeviceId    string `json:"device_id"`
	TokenString string `json:"token" binding:"required"`
	EventType   int    `json:"showType"`
	BaseQueryParam
}

type ReqAlertEventOriginData struct {
	DeviceId    string `json:"device_id"`
	TokenString string `json:"token" binding:"required"`
	TimeType    int    `json:"timeType"`
	DataType    int    `json:"dataType"`
	SelectTime  int    `json:"selectTime"`
	BaseQueryParam
}

type ResAlertDetail struct {
	Alert  GwAlert                `json:"alert"`
	Notify []ResAlertDetailNotify `json:"notify"`
	// NotifyType int `json:"notifyType"`
	// // GwNotify
}

type ResAlertDetailNotify struct {
	GwNotify  `xorm:"extends"`
	UserName  string `json:"userName"`
	UserPhone string `json:"phone"`
}

type ReqDeviceAlertDetail struct {
	DeviceId    string `json:"device_id"`
	TokenString string `json:"token" binding:"required"`
	BaseQueryParam
}

func (alert *Model) GetAlertList(params DeviceParams, ComId int) (int64, []GetDeviceList, error) {
	var (
		// gwAlerts []GwAlert
		devices      []GetDeviceList
		deviceCount  GwDevice
		xSession     *xorm.Session
		countSession *xorm.Session
	)
	fmt.Println(params)

	db, _ := utils.Connect()
	// db.ShowSQL(true)
	//params.ComId 这个是查询ID 查询公司列表ID
	//ComId 是公司自身ID
	//params.ComId == 0 查询全部
	//params.State查询的类型
	// companyId := ComId
	// queryState := params.State
	queryCompanySelectId := params.ComId

	if ComId != 1 && queryCompanySelectId != ComId {
		return 0, devices, fmt.Errorf("You have no permisson")
	}
	xSession = db.Table("gw_device").Select("gw_device.device_id,gw_device.descrip,gw_device.id,gw_device.address,gw_device.lng,gw_device.lat,gw_device.state,gw_device.company_id,gw_device.alert_id,gw_device.hearttime,gw_company.name as comname").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id")

	if ComId == 1 && queryCompanySelectId == 0 {
		xSession = xSession.Where("1=1")
		countSession = db.Where("1=1")
	} else {
		xSession = xSession.Where("gw_device.company_id = ?", queryCompanySelectId)
		countSession = db.Where("company_id = ?", queryCompanySelectId)
	}

	if params.State != "0" {
		xSession = xSession.And("gw_device.state = ?", params.State)
		countSession = countSession.And("state = ?", params.State)
	} else {
		xSession = xSession.And("gw_device.state IN ('10','20','30','70')")
		countSession = countSession.And("state IN ('10','20','30','70')")
	}

	total, _ := countSession.Count(&deviceCount)
	if err := xSession.Limit(params.Limit, params.Offset).Asc("hearttime").Find(&devices); err != nil {
		return total, devices, err
	}
	return total, devices, nil
}

func (alert *Model) GetDeviceAlertInfoByDeviceId(deviceId string) (DeviceAlertInfo, error) {
	var (
		deviceAlertInfo DeviceAlertInfo
		device          GwDevice
		company         GwCompany
		gwalert         GwAlert
		notify          []ResNotify
		//teles           []GwUser
		//该公司下的用户
		users []GwUser
	)
	db, _ := utils.Connect()

	if found, _ := db.Where("device_id = ?", deviceId).Get(&device); found == false {
		return deviceAlertInfo, fmt.Errorf("not found device")
	}

	if found, _ := db.Where("id = ?", device.CompanyId).Get(&company); found == false {
		return deviceAlertInfo, fmt.Errorf("not found company")
	}

	if found, _ := db.Where("id = ?", device.AlertId).Get(&gwalert); found == false {
		return deviceAlertInfo, fmt.Errorf("not found alert")
	}
	deviceAlertInfo.Address = device.Address
	deviceAlertInfo.Cola = gwalert.Cola
	deviceAlertInfo.State = device.State
	deviceAlertInfo.Value1 = company.Value1
	deviceAlertInfo.Value2 = company.Value2
	deviceAlertInfo.SendTime = gwalert.Sendtime
	deviceAlertInfo.DeviceId = device.DeviceId
	// err := db.Where("device_id = ? ", device.DeviceId).Desc("sendtime").Limit(3).Find(&nofify);err != nil {
	// 	return  deviceAlertInfo,err
	// }
	db.ShowSQL(true)
	//获取可以接受通知的联系人和联系方式
	err := db.Where("company_id = ?", device.CompanyId).And("auth_ids like '%4'").Find(&users)
	if err != nil {
		return deviceAlertInfo, err
	}
	//当没有接收通知用户的时候，err为"0"
	if len(users) < 1 {
		return deviceAlertInfo, fmt.Errorf("0")
	}
	notify = make([]ResNotify, len(users))
	for k, v := range users {
		//通知组
		var (
			gwNotify []GwNotify
		)
		notify[k].Name = v.RealName
		notify[k].Phone = v.Phone
		err := db.Where("user_id = ?", v.Id).And("device_id = ?", deviceId).Desc("id").Limit(3, 0).Find(&gwNotify)
		if err != nil {
			return deviceAlertInfo, err
		}
		notify[k].Data = gwNotify
	}
	deviceAlertInfo.Teles = users
	deviceAlertInfo.Notify = notify
	fmt.Println(notify)
	return deviceAlertInfo, nil
}

func (alert *Model) GetDevicePressureHistory(deviceId string, selectTime int) (DevicePressureHistory, error) {
	var (
		pressure []GwPressure
		timeMap  []map[string]string
		data     DevicePressureHistory
		company  GwCompany
	)

	db, _ := utils.Connect()

	switch selectTime {
	//全部
	case 5:
		time := make(map[string]string)
		time["startTime"] = "2018-06-03 00:00:00"
		time["endTime"] = "2050-06-03 00:00:00"
		timeMap = append(timeMap, time)
	//本周
	case 1:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 6 DAY), '%Y-%m-%d 23:59:59') AS 'endTime'")
	//上周
	case 2:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 1 WEEK), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 1), '%Y-%m-%d 23:59:59') AS 'endTime'")
	//本月
	case 3:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( CURDATE(), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(CURDATE()), '%Y-%m-%d 23:59:59') AS 'endTime'")
	//上月
	case 4:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(DATE_SUB(CURDATE(), INTERVAL 1 MONTH)), '%Y-%m-%d 23:59:59') AS 'endTime'")
	default:
		fmt.Println(0)
	}
	fmt.Println(timeMap[0]["startTime"])
	startTime, _ := utils.TimeToTimestamp(timeMap[0]["startTime"])
	endTime, _ := utils.TimeToTimestamp(timeMap[0]["endTime"])
	// db.ShowSQL(true)
	if err := db.Where("device_id = ?", deviceId).And("sendtime >= ?", startTime).And("sendtime <= ?", endTime).Find(&pressure); err != nil {
		return data, err
	}
	data.List = pressure
	if len(pressure) == 0 {
		return data, nil
	}

	db.Where("id = ?", pressure[0].CompanyId).Get(&company)
	data.Value1 = company.Value1
	data.Value2 = company.Value2
	return data, nil
}

func (alert *Model) GetDeviceOpenRecord(deviceId string) ([]GwAlert, error) {
	var (
		alerts []GwAlert
	)
	db, _ := utils.Connect()
	if err := db.Where("device_id = ?", deviceId).And("alert_type = ?", "20").And("cola != ?", "0").Find(&alerts); err != nil {
		return alerts, err
	}
	return alerts, nil
}

func (alert *Model) GetUserNotifyHistory(userId int, deviceId string) ([]GwNotify, error) {
	var (
		notify []GwNotify
	)
	db, _ := utils.Connect()
	if err := db.Where("user_id = ?", userId).And("device_id = ?", deviceId).Find(&notify); err != nil {
		return notify, err
	}
	return notify, nil
}

func (alert *Model) GetAlertHistory(params ReqAlertHistory, companyId int) (int64, []ResAlertForHistory, error) {
	var (
		resAlertForHistory []ResAlertForHistory
		xSession           *xorm.Session
		countSession       *xorm.Session
		countAlert         []ResAlertForHistory
	)

	db, _ := utils.Connect()
	db.ShowSQL(true)
	if companyId != 1 && params.CompanyId != companyId {
		return 0, resAlertForHistory, fmt.Errorf("You have no permisson!")
	}

	xSession = db.Table("gw_alert").Select("gw_alert.*,gw_device.address,gw_company.name as company_name").Join("INNER", "gw_device", "gw_alert.device_id = gw_device.device_id").Join("INNER", "gw_company", "gw_alert.company_id = gw_company.id")
	countSession = db.Table("gw_alert").Select("gw_alert.*,gw_device.address,gw_company.name as company_name").Join("INNER", "gw_device", "gw_alert.device_id = gw_device.device_id").Join("INNER", "gw_company", "gw_alert.company_id = gw_company.id")
	if params.CompanyId == 0 {
		xSession = xSession.Where("1=1")
		countSession = countSession.Where("1=1")
	} else {
		xSession = xSession.Where("gw_alert.company_id = ?", params.CompanyId)
		countSession = countSession.Where("gw_alert.company_id = ?", params.CompanyId)
	}

	if params.Addkeys != "" {
		//搜索地址
		xSession = xSession.And("gw_device.address like ?", "%"+params.Addkeys+"%")
		countSession = countSession.And("gw_device.address like ?", "%"+params.Addkeys+"%")
	}
	if len(params.DataPicker) > 0 {
		if params.DataPicker[0] != "" && params.DataPicker[1] != "" {
			xSession = xSession.And("gw_alert.sendtime >= ?", params.DataPicker[0]).And("gw_alert.sendtime <= ?", params.DataPicker[1])
			countSession = countSession.And("gw_alert.sendtime >= ?", params.DataPicker[0]).And("gw_alert.sendtime <= ?", params.DataPicker[1])
		}
	}
	xSession = xSession.And("gw_alert.cola != ?", "0")
	countSession = countSession.And("gw_alert.cola != ?", "0")
	switch params.AlertState {
	case 0:
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60','70')")
		countSession = countSession.And("gw_alert.alert_type IN ('20','30','60','70')")
	case 1:
		xSession = xSession.And("gw_alert.restoretime != ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
		countSession = countSession.And("gw_alert.restoretime != ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
	case 2:
		xSession = xSession.And("gw_alert.restoretime = ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
		countSession = countSession.And("gw_alert.restoretime = ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
	}
	count, _ := countSession.FindAndCount(&countAlert)
	fmt.Println(count)
	if err := xSession.Limit(params.Limit, params.Offset).Asc("gw_alert.restoretime").Desc("gw_alert.id").Find(&resAlertForHistory); err != nil {
		return count, resAlertForHistory, err
	}
	// fmt.Println(resAlertForHistory)
	if len(resAlertForHistory) > 0 {
		for k, v := range resAlertForHistory {
			var (
				notifystate GwNotify
			)
			notifycount, _ := db.Where("alert_id = ?", v.Id).And("state = ?", 1).And("type = ?", 1).Count(&notifystate)
			if notifycount == 0 {
				resAlertForHistory[k].NotifyStatus = 0
			} else {
				resAlertForHistory[k].NotifyStatus = 1
			}
		}
	}
	return count, resAlertForHistory, nil
}

func (user *Model) ExportDeviceAlertEvent(eventType int, deviceId string) ([]ResDeviceAlertEvent, error) {
	var (
		list     []ResDeviceAlertEvent
		xSession *xorm.Session
	)
	db, _ := utils.Connect()
	xSession = db.Table("gw_alert").Select("gw_alert.sendtime,gw_alert.id,gw_alert.alert_type,gw_alert.cola,gw_company.name as company_name").
		Join("INNER", "gw_company", "gw_alert.company_id = gw_company.id").Where("gw_alert.cola != ?", 0).And("gw_alert.cola != ? ", -1).And("gw_alert.device_id = ?", deviceId)
	if eventType == 1 {
		//全部
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60','70')")
	} else if eventType == 2 {
		//异常 70
		xSession = xSession.And("gw_alert.alert_type = ?", "70")
	} else if eventType == 3 {
		//告警
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60')")
	}
	if err := xSession.Find(&list); err != nil {
		return list, err
	}
	return list, nil
}

func (user *Model) GetDeviceAlertEvent(params ReqAlertEvent) (int64, []ResDeviceAlertEvent, error) {
	var (
		list         []ResDeviceAlertEvent
		xSession     *xorm.Session
		countSession *xorm.Session
		alert        GwAlert
	)
	db, _ := utils.Connect()
	// db.ShowSQL(true)
	xSession = db.Table("gw_alert").Select("gw_alert.sendtime,gw_alert.id,gw_alert.alert_type,gw_alert.cola,gw_company.name as company_name").
		Join("INNER", "gw_company", "gw_alert.company_id = gw_company.id").Where("gw_alert.cola != ?", 0).And("gw_alert.cola != ? ", -1).And("gw_alert.device_id = ?", params.DeviceId)
	countSession = db.Where("cola != ?", 0).And("cola != ? ", -1).And("device_id = ?", params.DeviceId)
	if params.EventType == 1 {
		//全部
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60','70')")
		countSession = countSession.And("alert_type IN ('20','30','60','70')")
	} else if params.EventType == 2 {
		//异常 70
		xSession = xSession.And("gw_alert.alert_type = ?", "70")
		countSession = countSession.And("alert_type = ?", "70")
	} else if params.EventType == 3 {
		//告警
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60')")
		countSession = countSession.And("alert_type IN ('20','30','60')")
	}
	count, _ := countSession.Count(&alert)
	err := xSession.Limit(params.Limit, params.Offset).Find(&list)
	if err != nil {
		return count, list, err
	}
	return count, list, nil
}

func (user *Model) ExportDeviceAlertOriginData(timeType int, selectTime int, dataType int, deviceId string) ([]GwAlert, error) {
	var (
		list      []GwAlert
		xSession  *xorm.Session
		startTime int64
		endTime   int64
		timeMap   []map[string]string
	)
	db, _ := utils.Connect()
	switch timeType {
	//全部
	case 1:
		time := make(map[string]string)
		time["startTime"] = "1700-06-03 00:00:00"
		time["endTime"] = "2050-06-03 00:00:00"
		timeMap = append(timeMap, time)
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//本周
	case 2:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 6 DAY), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//上周
	case 3:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 1 WEEK), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 1), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//本月
	case 4:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( CURDATE(), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(CURDATE()), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//上月
	case 5:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(DATE_SUB(CURDATE(), INTERVAL 1 MONTH)), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	case 6:
		//指定日期
		dayBegin := selectTime / 1000
		dayEnd := selectTime/1000 + 60*60*24 - 1

		startTime = int64(dayBegin)
		endTime = int64(dayEnd)
	default:
		fmt.Println(0)
	}

	xSession = db.Where("cola != ?", 0).And("cola != ?", -1).And("device_id = ?", deviceId).And("sendtime >= ?", startTime).And("sendtime <= ?", endTime)
	switch dataType {
	case 1:
		//全部
		xSession = xSession.And("alert_type IN ('20','30','40','60','70')")
	case 2:
		//心跳
		xSession = xSession.And("alert_type = ?", "40")
	case 3:
		//失联
		xSession = xSession.And("alert_type = ?", "70")
	case 4:
		//水压
		xSession = xSession.And("alert_type = ?", "60")
	case 5:
		//栓帽打开
		xSession = xSession.And("alert_type = ?", "20")
	}
	err := xSession.Find(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (user *Model) GetDeviceAlertOriginData(params ReqAlertEventOriginData) (int64, []GwAlert, error) {
	var (
		list         []GwAlert
		xSession     *xorm.Session
		countSession *xorm.Session
		alert        GwAlert
		startTime    int64
		endTime      int64
		timeMap      []map[string]string
	)
	db, _ := utils.Connect()
	db.ShowSQL(true)
	switch params.TimeType {
	//全部
	case 1:
		time := make(map[string]string)
		time["startTime"] = "2018-06-03 00:00:00"
		time["endTime"] = "2050-06-03 00:00:00"
		timeMap = append(timeMap, time)
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//本周
	case 2:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 6 DAY), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//上周
	case 3:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB( DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 1 WEEK), '%Y-%m-%d 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 1), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//本月
	case 4:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( CURDATE(), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(CURDATE()), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	//上月
	case 5:
		timeMap, _ = db.QueryString("SELECT DATE_FORMAT( DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00') AS 'startTime'," +
			"DATE_FORMAT( LAST_DAY(DATE_SUB(CURDATE(), INTERVAL 1 MONTH)), '%Y-%m-%d 23:59:59') AS 'endTime'")
		startTime, _ = utils.TimeToTimestamp(timeMap[0]["startTime"])
		endTime, _ = utils.TimeToTimestamp(timeMap[0]["endTime"])
	case 6:
		//指定日期
		dayBegin := params.SelectTime / 1000
		dayEnd := params.SelectTime/1000 + 60*60*24 - 1

		startTime = int64(dayBegin)
		endTime = int64(dayEnd)
	default:
		fmt.Println(0)
	}

	xSession = db.Where("cola != ?", 0).And("cola != ?", -1).And("device_id = ?", params.DeviceId).And("sendtime >= ?", startTime).And("sendtime <= ?", endTime)
	countSession = db.Where("cola != ?", 0).And("cola != ?", -1).And("device_id = ?", params.DeviceId).And("sendtime >= ?", startTime).And("sendtime <= ?", endTime)
	switch params.DataType {
	case 1:
		//全部
		xSession = xSession.And("alert_type IN ('20','30','40','60','70')")
		countSession = countSession.And("alert_type IN ('20','30','40','60','70')")
	case 2:
		//心跳
		xSession = xSession.And("alert_type = ?", "40")
		countSession = countSession.And("alert_type = ?", "40")
	case 3:
		//失联
		xSession = xSession.And("alert_type = ?", "70")
		countSession = countSession.And("alert_type = ?", "70")
	case 4:
		//水压
		xSession = xSession.And("alert_type = ?", "60")
		countSession = countSession.And("alert_type = ?", "60")
	case 5:
		//栓帽打开
		xSession = xSession.And("alert_type = ?", "20")
		countSession = countSession.And("alert_type = ?", "20")
	}
	count, _ := countSession.Count(&alert)
	err := xSession.Limit(params.Limit, params.Offset).Find(&list)
	if err != nil {
		return count, list, err
	}
	return count, list, nil
}

func (user *Model) GetDeviceAlertDetail(params ReqDeviceAlertDetail) (int64, []ResAlertDetail, error) {
	var (
		alerts         []GwAlert
		resAlertDetail []ResAlertDetail
		// notify         []ResAlertDetailNotify
		alert GwAlert
	)
	db, _ := utils.Connect()
	count, _ := db.Where("device_id = ?", params.DeviceId).And("cola != ?", 0).And("cola != ?", -1).And("alert_type IN ('20','30','60','70')").Count(&alert)

	if err := db.Where("device_id = ?", params.DeviceId).And("cola != ?", 0).And("cola != ?", -1).And("alert_type IN ('20','30','60','70')").Limit(params.Limit, params.Offset).Find(&alerts); err != nil {
		return count, resAlertDetail, err
	}

	if len(alerts) > 0 {
		for _, val := range alerts {
			var alertDetail ResAlertDetail
			alertDetail.Alert = val
			var notify []ResAlertDetailNotify
			if err := db.Table("gw_notify").Select("gw_notify.*, gw_user.name as user_name,gw_user.phone as user_phone").Join("INNER", "gw_user", "gw_notify.user_id = gw_user.id").Where("alert_id = ?", val.Id).Find(&notify); err != nil {
				return count, resAlertDetail, err
			}
			alertDetail.Notify = notify
			resAlertDetail = append(resAlertDetail, alertDetail)
		}
	}
	fmt.Println(resAlertDetail)
	return count, resAlertDetail, nil
}

func (alert *Model) ExportAlertTrace(selectCompayId int, addkeys string, alertState int, dataPicker []string) ([]ResAlertForHistory, error) {
	var (
		resAlertForHistory []ResAlertForHistory
		xSession           *xorm.Session
	)

	db, _ := utils.Connect()
	db.ShowSQL(true)

	xSession = db.Table("gw_alert").Select("gw_alert.*,gw_device.address,gw_company.name as company_name").Join("INNER", "gw_device", "gw_alert.device_id = gw_device.device_id").Join("INNER", "gw_company", "gw_alert.company_id = gw_company.id")

	if selectCompayId == 0 {
		xSession = xSession.Where("1=1")
	} else {
		xSession = xSession.Where("gw_alert.company_id = ?", selectCompayId)
	}

	if addkeys != "" {
		//搜索地址
		xSession = xSession.And("gw_device.address like ?", "%"+addkeys+"%")
	}
	if len(dataPicker) == 2 {
		if dataPicker[0] != "" && dataPicker[1] != "" {
			xSession = xSession.And("gw_alert.sendtime >= ?", dataPicker[0]).And("gw_alert.sendtime <= ?", dataPicker[1])
		}
	}
	xSession = xSession.And("cola != ?", "0")
	switch alertState {
	case 0:
		xSession = xSession.And("gw_alert.alert_type IN ('20','30','60','70')")
	case 1:
		xSession = xSession.And("gw_alert.restoretime != ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
	case 2:
		xSession = xSession.And("gw_alert.restoretime = ?", 0).And("gw_alert.alert_type IN ('20','30','60','70')")
	}

	if err := xSession.Asc("gw_alert.restoretime").Desc("gw_alert.id").Find(&resAlertForHistory); err != nil {
		return resAlertForHistory, err
	}

	if len(resAlertForHistory) > 0 {
		for k, v := range resAlertForHistory {
			var (
				notifystate GwNotify
			)
			notifycount, _ := db.Where("alert_id = ?", v.Id).And("state = ?", 1).And("type = ?", 1).Count(&notifystate)
			if notifycount == 0 {
				resAlertForHistory[k].NotifyStatus = 0
			} else {
				resAlertForHistory[k].NotifyStatus = 1
			}
		}
	}
	return resAlertForHistory, nil
}
