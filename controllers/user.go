package controllers

import (
	"gw/model"
	"gw/utils"

	"github.com/gin-gonic/gin"
)

type UserGetCompanySelectList struct {
	Token
}

type UserList struct {
	CompanyId   int    `json:"companyId"`
	TokenString string `json:"token" binding:"required"`
}

type AddOneUser struct {
	TokenString string `json:"token" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	CompanyId   int    `json:"company_id" binding:"required"`
	AuthList    string `json:"authority" binding:"required"`
	Status      int    `json:"status" binding:"required"`
}

type EditUser struct {
	TokenString string `json:"token" binding:"required"`
	Id          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Password    string `json:"password"`
	Phone       string `json:"phone" binding:"required"`
	CompanyId   int    `json:"company_id" binding:"required"`
	AuthList    string `json:"authority" binding:"required"`
	Status      int    `json:"status" binding:"required"`
}

type AddCom struct {
	Token
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Manager string `json:"manager" binding:"required"`
	Tel     string `json:"tel" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Value1  string `json:"value1" binding:"required"`
	Value2  string `json:"value2" binding:"required"`
}

type EidtCom struct {
	Id int `json:"id" binding:"required"`
	AddCom
}

type Token struct {
	TokenString string `json:"token" xorm:"-" binding:"required"`
}

//获取所有权限
func (tc *BaseController) UserGetAuthCheckList(c *gin.Context) {
	var (
		token Token
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "token 解析失败",
		})
		return
	}
	router := c.Request.RequestURI
	if _, _, err := tc.CheckAuth(token.TokenString, router, true); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if AuthList, err := m.GetAllAuth(); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": AuthList})
		return
	}
}

// 添加人员
func (tc *BaseController) UserAdd(c *gin.Context) {
	var (
		token AddOneUser
		user  model.GwUser
	)
	router := c.Request.RequestURI
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, _, authErr := tc.CheckAuth(token.TokenString, router, true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	user.CompanyId = token.CompanyId
	user.Name = token.Name
	user.Password = utils.String2md5(token.Password)
	user.Status = token.Status
	user.AuthIds = token.AuthList
	user.Phone = token.Phone
	if err := m.AddUser(user); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK"})
	return
}

//人员编辑
func (tc *BaseController) UserEdit(c *gin.Context) {
	var (
		token EditUser
		user  model.GwUser
	)
	router := c.Request.RequestURI
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	_, _, authErr := tc.CheckAuth(token.TokenString, router, true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	user.Id = token.Id
	user.CompanyId = token.CompanyId
	user.Name = token.Name
	if len(token.Password) > 0 {
		user.Password = utils.String2md5(token.Password)
	}
	user.Status = token.Status
	user.AuthIds = token.AuthList
	user.Phone = token.Phone
	if num, err := m.EidtUser(user); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": num})
		return
	}

}

func (tc *BaseController) UserGetCompanySelectList(c *gin.Context) {
	var (
		token UserGetCompanySelectList
	)
	router := c.Request.RequestURI
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "token 解析失败",
		})
		return
	}
	_, companyId, authErr := tc.CheckAuth(token.TokenString, router, true)
	if authErr != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    authErr.Error(),
		})
		return
	}

	list, err := m.UserGetCompanySelectList(companyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 0, "msg": "ok", "data": list})
	return
}

// UserList 获取用户列表
// 初始化的时候没有Select的CompanyID
func (tc *BaseController) UserList(c *gin.Context) {
	var (
		token UserList
	)

	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "token 解析失败",
		})
		return
	}
	selectCompanyId := token.CompanyId
	userId, companyId, authErr := tc.CheckAuth(token.TokenString, router, false)
	//fmt.Println(userId)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}

	users, err := m.GetUserList(userId, companyId, selectCompanyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": users})
	return
}

func (tc *BaseController) UserInfo(c *gin.Context) {
	var (
		token Token
	)
	// authToken := c.Request.Header.Get("token")
	// fmt.Println(authToken)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	userId, _, err := tc.CheckAuth(token.TokenString, "", true)
	user, err := m.GetUserInfoByUserid(userId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": user})
	return
}

func (tc *BaseController) ComAdd(c *gin.Context) {
	var (
		token AddCom
		com   model.GwCompany
	)

	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, router, false); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	com.Name = token.Name
	com.Address = token.Address
	com.Email = token.Email
	com.Tel = token.Tel
	com.Value1 = token.Value1
	com.Value2 = token.Value2
	com.Manager = token.Manager
	if err := m.AddCom(com); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK"})
	return
}

func (tc *BaseController) ComEdit(c *gin.Context) {
	var (
		token EidtCom
		com   model.GwCompany
	)

	router := c.Request.RequestURI

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	if _, _, err := tc.CheckAuth(token.TokenString, router, false); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	com.Name = token.Name
	com.Address = token.Address
	com.Email = token.Email
	com.Tel = token.Tel
	com.Value1 = token.Value1
	com.Value2 = token.Value2
	com.Manager = token.Manager
	com.Id = token.Id
	num, err := m.EidtCom(com)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": num})
	return
}
