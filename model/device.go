package model

import (
	"fmt"
	"gw/utils"
	"strconv"

	"math"

	"github.com/go-xorm/xorm"
)

type DeviceList struct {
	// MinLatitude  float64 `json:"minLatitude"`
	// MaxLatitude  float64 `json:"maxLatitude"`
	// MinLongitude float64 `json:"minLongitude"`
	// MaxLongitude float64 `json:"maxLongitude"`
	CompanyId   int    `json:"companyId"`
	Status      int    `json: "status"`
	TokenString string `json:"token" binding:"required"`
	AddKeys     string `json:"addkeys"`
	OnlineState int    `json:"online_state"`
	BaseQueryParam
}

type DeviceMapList struct {
	MinLatitude  float64 `json:"minLatitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
	MinLongitude float64 `json:"minLongitude"`
	MaxLongitude float64 `json:"maxLongitude"`
	CompanyId    int     `json:"companyId"`
	Status       int     `json: "status"`
	TokenString  string  `json:"token" binding:"required"`
	AddKeys      string  `json:"addkeys"`
	OnlineState  int     `json:"online_state"`
}

type ReqDeviceUnalert struct {
	MinLatitude  float64 `json:"minLatitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
	MinLongitude float64 `json:"minLongitude"`
	MaxLongitude float64 `json:"maxLongitude"`
	CompanyId    int     `json:"companyId"`
	TokenString  string  `json:"token" binding:"required"`
}

type ExportDeviceQuery struct {
	TokenString string `form:"token"`
	CompanyId   int    `form:"companyId"`
	AddKeys     string `form:"addkeys"`
	OnlineState int    `form:"online_state"`
}

type ResDeviceMonitoring struct {
	Device           []GwDevice `json:"alertList"`
	DeviceTotal      int64      `json:"deviceTotal"`
	AlertDeviceCount int64      `json:"alertTotal"`
	OpenCount        int64      `json:"opc"`
	SignalCount      int64      `json:"sigc"`
	PressureCount    int64      `json:"prec"`
	DownCount        int64      `json:"downc"`
	PowerCount       int64      `json:"powerc"`
	LossCount        int64      `json:"lossc"`
}

func (user *Model) DeviceList(selectCompanyId int, companyId int, params DeviceList) (int64, []GetDeviceList, error) {
	var (
		devices []GetDeviceList
		device  GwDevice
		// geoDevice    []GetDeviceList
		xSession     *xorm.Session
		countSession *xorm.Session
	)

	queryCompanyId := params.CompanyId

	db, _ := utils.Connect()
	db.ShowSQL(true)
	if companyId == 1 && queryCompanyId == 0 {
		xSession = db.Where("1=1")
		countSession = db.Where("1=1")
	} else {
		xSession = db.Where("gw_device.company_id = ?", queryCompanyId)
		countSession = db.Where("company_id = ?", queryCompanyId)
	}

	if params.Status == 1 {
		xSession = xSession.And("gw_device.status = ?", 0)
		countSession = countSession.And("status = ?", 0)
	}
	if params.Status == 2 {
		xSession = xSession.And("gw_device.status = ?", 1)
		countSession = countSession.And("status = ?", 1)
	}
	if params.AddKeys != "" {
		xSession = xSession.And("gw_device.address like ?", "%"+params.AddKeys+"%")
		countSession = countSession.And("gw_device.address like ?", "%"+params.AddKeys+"%")
	}
	if params.OnlineState == 1 {
		xSession = xSession.And("gw_device.state != ?", "70").And("gw_device.state != ?", "80").And("gw_device.status = ?", 1)
		countSession = countSession.And("gw_device.state != ?", "70").And("gw_device.state != ?", "80").And("gw_device.status = ?", 1)
	} else if params.OnlineState == 2 {
		xSession = xSession.And("gw_device.state = ?", "70").Or("gw_device.state = ?", "80").And("gw_device.status = ?", 1)
		countSession = countSession.And("gw_device.state = ?", "70").Or("gw_device.state = ?", "80").And("gw_device.status = ?", 1)
	}

	xSession = xSession.Table("gw_device").Select("gw_device.*,gw_company.name as comname,gw_company.manager,gw_company.tel").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id")
		// Where("gw_device.company_id = ?", queryCompanyId)

	num, _ := countSession.Count(&device)
	if params.Order == "desc" {
		xSession = xSession.Desc(params.Sort).Asc("id")
	}

	if params.Limit > 0 {
		err := xSession.Limit(params.Limit, params.Offset).Find(&devices)
		if err != nil {
			return 0, devices, err
		}

		// if params.Limit == 10000 {
		// 	for _, v := range devices {
		// 		////myLatitude, myLongitude, minLatitude, maxLatitude, minLongitude, maxLongitude
		// 		myLatitude, _ := strconv.ParseFloat(v.Lat, 64)
		// 		myLongitude, _ := strconv.ParseFloat(v.Lng, 64)
		// 		if isInArea(myLatitude, myLongitude, params.MinLatitude, params.MaxLatitude, params.MinLongitude, params.MaxLongitude) {
		// 			geoDevice = append(geoDevice, v)
		// 		}
		// 	}
		// 	fmt.Println(len(geoDevice))
		// 	return int64(len(geoDevice)), geoDevice, nil
		// }
	} else {
		err := xSession.Find(&devices)
		if err != nil {
			return 0, devices, err
		}

	}

	return num, devices, nil
}

func (user *Model) DeviceMapList(selectCompanyId int, companyId int, params DeviceMapList) ([]GetDeviceList, error) {
	var (
		devices []GetDeviceList
		// device       GwDevice
		// geoDevice []GetDeviceList
		xSession *xorm.Session
	)

	queryCompanyId := params.CompanyId

	db, _ := utils.Connect()
	db.ShowSQL(true)
	if companyId == 1 && queryCompanyId == 0 {
		xSession = db.Where("1=1")
	} else {
		xSession = db.Where("gw_device.company_id = ?", queryCompanyId)
	}

	if params.Status == 1 {
		xSession = xSession.And("gw_device.status = ?", 0)
	}
	if params.Status == 2 {
		xSession = xSession.And("gw_device.status = ?", 1)
	}
	if params.AddKeys != "" {
		xSession = xSession.And("gw_device.address like ?", "%"+params.AddKeys+"%")
	}
	if params.OnlineState == 1 {
		xSession = xSession.And("gw_device.state != ?", "70").And("gw_device.state != ?", "80").And("gw_device.status = ?", 1)
	} else if params.OnlineState == 2 {
		xSession = xSession.And("gw_device.state = ?", "70").Or("gw_device.state = ?", "80").And("gw_device.status = ?", 1)
	}

	err := xSession.Table("gw_device").Select("gw_device.*,gw_company.name as comname,gw_company.manager,gw_company.tel").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").
		And("gw_device.lng >= ?", params.MinLongitude).
		And("gw_device.lng <= ?", params.MaxLongitude).
		And("gw_device.lat >= ?", params.MinLatitude).
		And("gw_device.lat <= ?", params.MaxLatitude).
		Find(&devices)
	if err != nil {
		return devices, err
	}
	return devices, nil
}

func (user *Model) AddOneDevice(device GwDevice) error {
	db, _ := utils.Connect()
	if _, err := db.Insert(&device); err != nil {
		return err
	}
	return nil
}

func (user *Model) EidtDevicae(device GwDevice) (int64, error) {
	db, _ := utils.Connect()
	affected, err := db.Id(device.Id).Update(device)
	if err != nil {
		return affected, err
	}
	return affected, nil
}

func (user *Model) ImportDevice(devices []GwDevice, importType int, companyId int) (int64, error) {
	var (
		dbDevices []GwDevice
		dbDevice  GwDevice
	)
	db, _ := utils.Connect()
	//如果查询出来的记录为空那么是新建的公司
	//可以不用判断直接插入数据
	err := db.Where("company_id = ?", companyId).Find(&dbDevices)
	if err != nil {
		return 0, err
	}
	if len(dbDevices) == 0 {
		count, _ := db.Insert(&devices)
		return count, nil
	}
	if importType == 1 {
		var (
			num int64
		)
		//数据库记录和传值做对比
		for _, v := range dbDevices {
			found := func(deviceId string) bool {
				for _, s := range devices {
					if deviceId == s.DeviceId {
						return true
					}
				}
				return false
			}(v.DeviceId)
			//不在传值里面则删除这条记录
			if found == false {
				db.Where("device_id = ?", v.DeviceId).Delete(dbDevice)
			}
		}
		//传值和数据库做对比
		for _, v := range devices {
			found := func(deviceId string) bool {
				for _, s := range dbDevices {
					if deviceId == s.DeviceId {
						return true
					}
				}
				return false
			}(v.DeviceId)
			//不在数据库中则插入
			if found == false {
				num, _ = db.Insert(&v)
			}
		}
		return num, nil
	} else if importType == 2 {
		db.Where("company_id = ?", companyId).Delete(dbDevice)
		num, err := db.Insert(&devices)
		if err != nil {
			return num, err
		}
		return num, nil
	}
	return 0, fmt.Errorf("type error")
}

func (user *Model) GetdeviceStateList(selectCompanyId int, companyId int, params DeviceList) (int64, []GetDeviceList, error) {
	var (
		devices      []GetDeviceList
		device       GwDevice
		xSession     *xorm.Session
		countSession *xorm.Session
	)

	queryCompanyId := params.CompanyId

	db, _ := utils.Connect()
	db.ShowSQL(true)
	if companyId == 1 && queryCompanyId == 0 {
		xSession = db.Where("1=1")
		countSession = db.Where("1=1")
	} else {
		xSession = db.Where("gw_device.company_id = ?", queryCompanyId)
		countSession = db.Where("company_id = ?", queryCompanyId)
	}

	// if params.Status != 2 {
	// 	xSession = xSession.And("gw_device.status = ?", params.Status)
	// 	countSession = countSession.And("status = ?", params.Status)
	// }
	if params.AddKeys != "" {
		xSession = xSession.And("gw_device.address like ?", "%"+params.AddKeys+"%")
		countSession = countSession.And("gw_device.address like ?", "%"+params.AddKeys+"%")
	}
	if params.OnlineState == 1 {
		xSession = xSession.And("gw_device.state != ?", "70")
		countSession = countSession.And("gw_device.state != ?", "70")
	} else if params.OnlineState == 2 {
		xSession = xSession.And("gw_device.state = ?", "70")
		countSession = countSession.And("gw_device.state = ?", "70")
	}

	xSession = xSession.Table("gw_device").Select("gw_device.*,gw_company.name as comname,gw_company.manager,gw_company.tel").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id")
		// Where("gw_device.company_id = ?", queryCompanyId)

	num, _ := countSession.Count(&device)
	if params.Order == "desc" {
		xSession = xSession.Desc(params.Sort).Asc("id")
	}

	err := xSession.Limit(params.Limit, params.Offset).Find(&devices)

	if err != nil {
		return 0, devices, err
	}
	return num, devices, nil
}

func (user *Model) GetExportDeviceList(companyId int, selectCompanyId string, onlineState string, addkeys string) ([]GetDeviceList, error) {
	var (
		devices  []GetDeviceList
		xSession *xorm.Session
	)

	queryCompanyId, _ := strconv.Atoi(selectCompanyId)

	db, _ := utils.Connect()
	db.ShowSQL(true)
	if companyId == 1 && queryCompanyId == 0 {
		xSession = db.Where("1=1")
	} else {
		xSession = db.Where("gw_device.company_id = ?", queryCompanyId)
	}

	if addkeys != "" {
		xSession = xSession.And("gw_device.address like ?", "%"+addkeys+"%")
	}
	if onlineState == "1" {
		xSession = xSession.And("gw_device.state != ?", "70")
	} else if onlineState == "2" {
		xSession = xSession.And("gw_device.state = ?", "70")
	}

	err := xSession.Table("gw_device").Select("gw_device.*,gw_company.name as comname,gw_company.manager,gw_company.tel").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").Find(&devices)

	if err != nil {
		return devices, err
	}
	return devices, nil
}

func isInRange(point float64, left float64, right float64) bool {
	if point >= math.Min(left, right) && point <= math.Max(left, right) {
		return true
	}
	return false
}

//myLatitude, myLongitude, minLatitude, maxLatitude, minLongitude, maxLongitude
func isInArea(latitue float64, longitude float64, areaLatitude1 float64, areaLatitude2 float64, areaLongitude1 float64, areaLongitude2 float64) bool {
	if isInRange(latitue, areaLatitude1, areaLatitude2) {
		if areaLongitude1*areaLongitude2 > 0 {
			if isInRange(longitude, areaLongitude1, areaLongitude2) {
				return true
			}
			return false
		} else {
			if math.Abs(areaLongitude1)+math.Abs(areaLongitude2) < 180 {
				if isInRange(longitude, areaLongitude1, areaLongitude2) {
					return true
				}
				return false
			} else {
				left := math.Max(areaLongitude1, areaLongitude2)
				right := math.Max(areaLongitude1, areaLongitude2)
				if isInRange(longitude, left, 180) || isInRange(longitude, right, -180) {
					return true
				}
				return false
			}
		}
	} else {
		return false
	}
}

func (user *Model) GetDeviceMonitoring(companyId int, selectCompanyId int) (ResDeviceMonitoring, error) {
	var (
		deviceCount      int64
		device           GwDevice
		alertDeviceCount int64
		alertDevices     []GwDevice
		xsession         *xorm.Session
		openCount        int64
		signalCount      int64
		pressureCount    int64
		downCount        int64
		powerCount       int64
		lossCount        int64
		resData          ResDeviceMonitoring
	)
	queryCompanyId := selectCompanyId
	if companyId != 1 {
		queryCompanyId = companyId
	}
	db, _ := utils.Connect()

	// if queryCompanyId == 0 {
	// 	xsession = db.Where("1=1")
	// } else {
	// 	xsession = db.Where("company_id = ?", queryCompanyId)
	// }

	if queryCompanyId == 0 {
		alertDeviceCount, _ = db.Where("state IN ('10','20','30','70')").Count(&device)
		xsession = db.Where("state IN ('10','20','30','70')")
		deviceCount, _ = db.Count(&device)
		openCount, _ = db.Where("state = ?", "20").Count(&device)
		pressureCount, _ = db.Where("state = ?", "10").Count(&device)
		downCount, _ = db.Where("state = ?", "30").Count(&device)
		lossCount, _ = db.Where("state = ?", "70").Count(&device)
		signalCount, _ = db.Where("`signal` < ?", "18").Count(&device)
		powerCount, _ = db.Where("battery < ?", 15).Count(&device)

	} else {
		deviceCount, _ = db.Where("company_id = ?", queryCompanyId).Count(&device)
		alertDeviceCount, _ = db.Where("company_id = ?", queryCompanyId).And("state IN ('10','20','30','70')").Count(&device)
		xsession = db.Where("company_id = ?", queryCompanyId).And("state IN ('10','20','30','70')")
		openCount, _ = db.Where("state = ?", "20").And("company_id = ?", queryCompanyId).Count(&device)
		pressureCount, _ = db.Where("state = ?", "10").And("company_id = ?", queryCompanyId).Count(&device)
		downCount, _ = db.Where("state = ?", "30").And("company_id = ?", queryCompanyId).Count(&device)
		lossCount, _ = db.Where("state = ?", "70").And("company_id = ?", queryCompanyId).Count(&device)
		signalCount, _ = db.Where("`signal` < ?", "15").And("company_id = ?", queryCompanyId).Count(&device)
		powerCount, _ = db.Where("battery < ?", 15).And("company_id = ?", queryCompanyId).Count(&device)
	}
	if err := xsession.Desc("hearttime").Find(&alertDevices); err != nil {
		return resData, err
	}
	resData.Device = alertDevices
	resData.DeviceTotal = deviceCount
	resData.AlertDeviceCount = alertDeviceCount
	resData.DownCount = downCount
	resData.LossCount = lossCount
	resData.OpenCount = openCount
	resData.PowerCount = powerCount
	resData.PressureCount = pressureCount
	resData.SignalCount = signalCount
	return resData, nil

}

func (user *Model) GetUnalertDevice(params ReqDeviceUnalert, companyId int) ([]GetDeviceList, error) {
	var (
		devices         []GetDeviceList
		selectCompanyId int
		xSession        *xorm.Session
	)
	db, _ := utils.Connect()
	selectCompanyId = params.CompanyId
	if companyId != 1 {
		selectCompanyId = companyId
	}
	db.ShowSQL(true)
	if selectCompanyId == 0 {
		xSession = db.Where("1=1")
	} else {
		xSession = db.Where("company_id = ?", selectCompanyId)
	}
	xSession = xSession.NotIn("state", "10", "20", "30", "70")

	err := xSession.Table("gw_device").Select("gw_device.*,gw_company.name as comname,gw_company.manager,gw_company.tel").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").
		And("gw_device.lng >= ?", params.MinLongitude).
		And("gw_device.lng <= ?", params.MaxLongitude).
		And("gw_device.lat >= ?", params.MinLatitude).
		And("gw_device.lat <= ?", params.MaxLatitude).
		Find(&devices)
	if err != nil {
		return devices, err
	}
	return devices, nil
}
