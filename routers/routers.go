package routers

import (
	"gw/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRouters(r *gin.Engine, tc *controllers.BaseController) {
	r.POST("/dologin", tc.Login)
	r.POST("/changePsd", tc.ChangePsd)
	r.POST("/home", tc.Home)
	r.POST("/userlist", tc.UserList)
	r.POST("/userInfo", tc.UserInfo)
	r.POST("/companyList", tc.UserGetCompanySelectList)
	r.POST("/authList", tc.UserGetAuthCheckList)
	r.POST("/useradd", tc.UserAdd)
	r.POST("/useredit", tc.UserEdit)
	r.POST("/comadd", tc.ComAdd)
	r.POST("/comedit", tc.ComEdit)
	r.POST("/devicelist", tc.DeviceList)
	r.POST("/deviceMapList", tc.DeviceMapList)
	r.POST("/deviceadd", tc.DeviceAdd)
	r.POST("/deviceedit", tc.DeviceEdit)
	r.POST("/deviceimport", tc.DeviceImport)
	// r.POST("/devicestatelist", tc.DeviceStateList)
	r.POST("/alertlist", tc.AlertList)
	r.POST("/alertInfo", tc.AlertInfo)
	r.POST("/devicePressurehistory", tc.DevicePressureHistory)
	r.POST("/deviceOpenhistory", tc.DeviceOpenHistory)
	r.POST("/userNotifyHistory", tc.UserNotifyHistory)
	r.POST("/alertTrace", tc.AlertHistory)
	r.POST("/deviceAlertEvent", tc.DeviceAlertEvent)
	r.POST("/deviceAlertOriginData", tc.DeviceAlertOriginData)
	r.POST("/deviceAlertDetail", tc.DeviceAlertDetail)
	r.POST("/deviceexport", tc.DeviceExport)
	r.POST("/deviceEventExport", tc.DeviceEventExport)
	r.POST("/exportAlertTrace", tc.ExportAlertTrace)
}
