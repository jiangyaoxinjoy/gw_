package model

import (
	"fmt"
	"gw/utils"
)

type DeviceList struct {
	CompanyId   int    `json:"companyId"`
	Status      int    `json: "status"`
	TokenString string `json:"token" binding:"required"`
	BaseQueryParam
}

func (user *Model) DeviceList(selectCompanyId int, companyId int, params DeviceList) (int64, []GetDeviceList, error) {
	var (
		devices []GetDeviceList
		device  GwDevice
	)

	queryCompanyId := companyId

	if selectCompanyId != 0 && companyId == 1 {
		queryCompanyId = selectCompanyId
	}

	db, _ := utils.Connect()

	numqs := db.Where("company_id = ?", queryCompanyId)

	qs := db.Table("gw_device").Select("gw_device.*,gw_company.name as comname").
		Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").
		Where("gw_device.company_id = ?", queryCompanyId)

	if params.Status != 2 {
		qs = qs.And("gw_device.status = ?", params.Status)
		numqs = numqs.And("status = ?", params.Status)
	}
	num, _ := numqs.Count(&device)
	if params.Order == "desc" {
		qs = qs.Desc(params.Sort).Asc("id")
	}
	err := qs.Limit(params.Limit, params.Offset).Find(&devices)
	// if params.Order == "asc" {

	// 	// err = db.Table("gw_device").Select("gw_device.*,gw_company.name as comname").
	// 	// 	Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").
	// 	// 	Where("gw_device.company_id = ?", queryCompanyId).Limit(params.Limit, params.Offset).Find(&devices)
	// }
	// if params.Order == "desc" {
	// 	// err = db.Table("gw_device").Select("gw_device.*,gw_company.name as comname").
	// 	// 	Join("INNER", "gw_company", "gw_device.company_id = gw_company.id").
	// 	// 	Where("gw_device.company_id = ?", queryCompanyId).Desc(params.Sort).Limit(params.Limit, params.Offset).Find(&devices)
	// }

	if err != nil {
		return 0, devices, err
	}
	return num, devices, nil
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
