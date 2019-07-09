package controllers

import (
	"gw/config"
	"gw/model"

	"fmt"
	"os"
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

type ReqMonitoring struct {
	CompanyId   int    `json:"companyId"`
	TokenString string `json:"token" binding:"required"`
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

func (tc *BaseController) DeviceMapList(c *gin.Context) {
	var (
		token           model.DeviceMapList
		selectCompanyId int
	)

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	selectCompanyId = token.CompanyId

	_, companyId, authErr := tc.CheckAuth(token.TokenString, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	list, err := m.DeviceMapList(selectCompanyId, companyId, token)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": list})
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

func (tc *BaseController) DeviceStateList(c *gin.Context) {
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
	len, list, err := m.GetdeviceStateList(selectCompanyId, companyId, token)
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

func (tc *BaseController) DeviceExport(c *gin.Context) {
	router := c.Request.RequestURI
	selectCompanyId := c.PostForm("companyId")
	onlineState := c.PostForm("online_state")
	token := c.PostForm("token")
	addkeys := c.PostForm("addkeys")

	_, companyId, authErr := tc.CheckAuth(token, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	list, err := m.GetExportDeviceList(companyId, selectCompanyId, onlineState, addkeys)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	for key, val := range list {
		fmt.Println(val)
		i := key + 1
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i), val.CompanyName)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), val.DeviceId)
		if val.Status == 0 {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "未安装")
		} else {
			if val.State == "70" {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "离线")
			} else {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), "在线")
			}
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i), val.Signal)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i), val.Address)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%v", i), val.Manager)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%v", i), val.Tel)
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
	//extraHeaders["Content-Description"] = "File Transfer"
	//extraHeaders["Content-Transfer-Encoding"] = "binary"
	extraHeaders["Expires"] = "0"
	extraHeaders["Cache-Control"] = "must-revalidate"
	extraHeaders["Pragma"] = "public"

	// c.Header("Content-Disposition", "attachment; filename="+fileName)
	// c.Header("Content-Description", "File Transfer")
	// c.Header("Content-Type", "application/octet-stream")
	// c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Expires", "0")
	// c.Header("Cache-Control", "must-revalidate")
	// c.Header("Pragma", "public")
	nf, _ := os.Open(config.ExportFolder + fileName)
	defer nf.Close()
	fi, _ := nf.Stat()
	//fi.Size()

	c.DataFromReader(200, fi.Size(), "application/octet-stream", nf, extraHeaders)

	// c.JSON(200, gin.H{"status": 0, "msg": "OK"})
	return
}

func (tc *BaseController) DeviceMonitoring(c *gin.Context) {
	var (
		query ReqMonitoring
	)

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	// token := c.Request.Header.Get("token")
	_, companyId, authErr := tc.CheckAuth(query.TokenString, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}

	data, err := m.GetDeviceMonitoring(companyId, query.CompanyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) ShowUnalertDeivce(c *gin.Context) {
	var (
		params model.ReqDeviceUnalert
	)
	token := c.Request.Header.Get("token")
	_, companyId, authErr := tc.CheckAuth(token, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	data, err := m.GetUnalertDevice(params, companyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}
