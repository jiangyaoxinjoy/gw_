package model

import (
	"fmt"
	"gw/utils"
	// "github.com/go-xorm/xorm"
	// "errors"
)

// LoginReq 登录请求参数类
type LoginReq struct {
	Name        string `json:"userName" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CompanyName string `json:"company" binding:"required"`
}

// LoginCheck 登录验证
func (login *Model) LoginCheck(req LoginReq) (GwUser, error) {
	var (
		com  GwCompany
		user GwUser
	)
	db, _ := utils.Connect()
	comfound, _ := db.Where("name = ?", req.CompanyName).Get(&com)
	if comfound == false {
		return user, fmt.Errorf("Can not find company!")
	}
	//db.ShowSQL(true)
	userfound, _ := db.Where("name = ?", req.Name).And("password = ?", req.Password).And("company_id = ?", com.Id).Get(&user)
	if userfound == false {
		return user, fmt.Errorf("Can not find user!")
	}
	if user.Status != 1 {
		return user, fmt.Errorf("用户已被禁用!")
	}
	//fmt.Println(user.Id)
	return user, nil
}
