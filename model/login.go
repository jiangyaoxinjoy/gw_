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

type ReqChangePsd struct {
	TokenString string `json:"token" binding:"required"`
	OriginPsd   string `json:"originPwd"`
	NewPsd      string `json:"newPwd"`
	RepeatPsd   string `json:"againPwd"`
	UserId      int    `json:"userId" binding:"required"`
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
		return user, fmt.Errorf("用户名或密码错误!")
	}
	if user.Status != 1 {
		return user, fmt.Errorf("用户已被禁用!")
	}
	//fmt.Println(user.Id)
	return user, nil
}

func (LoginReq *Model) ChangePassword(token ReqChangePsd, companyId int) error {
	var (
		user GwUser
	)
	db, _ := utils.Connect()
	db.ShowSQL(true)
	token.OriginPsd = utils.String2md5(token.OriginPsd)
	if found, _ := db.Where("id = ?", token.UserId).And("password = ?", token.OriginPsd).And("company_id = ?", companyId).Get(&user); found == false {
		return fmt.Errorf("原密码错误!")
	}
	user.Password = utils.String2md5(token.NewPsd)
	if _, err := db.Id(user.Id).Update(&user); err != nil {
		return err
	}
	return nil
}
