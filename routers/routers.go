package routers

import (
	"gw/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRouters(r *gin.Engine, tc *controllers.BaseController) {
	r.POST("/login", tc.Login)
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
	r.POST("/deviceadd", tc.DeviceAdd)
	r.POST("/deviceedit", tc.DeviceEdit)
	r.POST("/deviceimport", tc.DeviceImport)
}
