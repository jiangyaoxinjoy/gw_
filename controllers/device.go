package controllers

import (
	"fmt"
	"gw/config"
	"gw/model"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type ReqDevice struct {
	Token
	model.GwDevice
}

type UpLoadFile struct {
	TokenString string `form:"token" binding:"required"`
	CompanyId   int    `form:"companyId" binding:"required"`
	ImportType  int    `form:"importType" binding:"required"`
}

func (tc *BaseController) DeviceList(c *gin.Context) {
	var (
		token           model.DeviceList
		selectCompanyId int
	)
	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	selectCompanyId = token.CompanyId

	_, companyId, authErr := tc.CheckAuth(token.TokenString, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	len, list, err := m.DeviceList(selectCompanyId, companyId, token)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data := struct {
		List  []model.GetDeviceList `json:"list"`
		Count int64                 `json:"count"`
	}{
		List:  list,
		Count: len,
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) DeviceAdd(c *gin.Context) {
	var (
		token  ReqDevice
		device model.GwDevice
	)
	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	_, _, authErr := tc.CheckAuth(token.TokenString, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	device.Address = token.Address
	device.Lng = token.Lng
	device.Lat = token.Lat
	device.DeviceId = token.DeviceId
	device.CompanyId = token.CompanyId
	if err := m.AddOneDevice(device); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK"})
	return
}

func (tc *BaseController) DeviceEdit(c *gin.Context) {
	var (
		token  ReqDevice
		device model.GwDevice
	)
	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if token.Id == 0 {
		c.JSON(200, gin.H{"status": -1, "msg": "设备ID无效"})
		return
	}
	_, _, authErr := tc.CheckAuth(token.TokenString, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	device.Address = token.Address
	device.Lng = token.Lng
	device.Lat = token.Lat
	device.DeviceId = token.DeviceId
	device.CompanyId = token.CompanyId
	device.Id = token.Id
	num, err := m.EidtDevicae(device)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": num})
	return
}

func (tc *BaseController) DeviceImport(c *gin.Context) {
	var (
		devices []model.GwDevice
		token   UpLoadFile
		device  model.GwDevice
	)
	router := c.Request.RequestURI
	if err := c.ShouldBind(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	_, _, authErr := tc.CheckAuth(token.TokenString, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}

	file, _ := c.FormFile("file")
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".xlsx"
	if err := c.SaveUploadedFile(file, config.Folder+fileName); err != nil {
		fmt.Println(err)
	}
	// fmt.Println(file.Filename)

	f, err := excelize.OpenFile(config.Folder + fileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	rows, _ := f.GetRows("Sheet1")
	if len(rows) < 1 {
		c.JSON(200, gin.H{"status": -1, "msg": fmt.Errorf("Can not find any devices!")})
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			if colCell != "" {
				device.DeviceId = colCell
				device.CompanyId = token.CompanyId
				devices = append(devices, device)
			}
		}
	}
	num, errM := m.ImportDevice(devices, token.ImportType, token.CompanyId)
	if errM != nil {
		c.JSON(200, gin.H{"status": -1, "msg": errM.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": num})
	return
}
