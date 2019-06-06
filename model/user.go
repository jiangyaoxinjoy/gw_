package model

import (
	"fmt"
	"gw/utils"
	"strings"
)

func (user *Model) CheckAuthByUserId(userId float64, r string) error {
	var (
		gwuser  GwUser
		authsub []GwAuthSub
	)
	db, _ := utils.Connect()
	//db.ShowSQL(true)
	found, _ := db.Where("id = ?", int(userId)).Get(&gwuser)
	//fmt.Println(gwuser)
	if found == false {
		return fmt.Errorf("can not found user")
	}
	auth := strings.Split(gwuser.AuthIds, ",")
	db.In("auth_id", auth).Find(&authsub)
	if len(authsub) > 0 {
		for _, v := range authsub {
			if r == v.Node {
				return nil
			}
		}
	}
	return fmt.Errorf("no auth!")
}

type GetUserList struct {
	CompanyName string   `xorm:"comname" json:"company"`
	AuthName    []string `json:"authName"`
	GwUser      `xorm:"extends"`
}

type GetDeviceList struct {
	CompanyName string `xorm:"comname" json:"company"`
	GwDevice    `xorm:"extends"`
}

func (user *Model) GetUserList(userId int, companyId int, selectCompanyId int) ([]GetUserList, error) {
	// fmt.Println(userId)
	var (
		users []GetUserList
		// auth  []GwAuthority
	)

	queryCompanyId := companyId

	if selectCompanyId != 0 && companyId == 1 {
		queryCompanyId = selectCompanyId
	}
	db, _ := utils.Connect()

	err := db.Table("gw_user").Select("gw_user.id,gw_user.name,gw_user.company_id,gw_user.phone,gw_user.status,gw_user.auth_ids,gw_company.name as comname").
		Join("INNER", "gw_company", "gw_user.company_id = gw_company.id").
		Where("gw_user.company_id = ?", queryCompanyId).Find(&users)
	// db.ShowSQL(true)
	if err != nil {
		return users, err
	}
	for k, v := range users {
		var (
			auth []GwAuthority
		)
		db.In("id", strings.Split(v.AuthIds, ",")).Find(&auth)
		for i := range auth {
			users[k].AuthName = append(users[k].AuthName, auth[i].Name)
		}
	}
	return users, nil
}

func (user *Model) UserGetCompanySelectList(companyId int) ([]GwCompany, error) {
	var (
		company []GwCompany
	)
	db, _ := utils.Connect()
	if companyId == 1 {
		err := db.Find(&company)
		if err != nil {
			return company, err
		}
		return company, nil
	}

	err := db.Where("id = ?", companyId).Find(&company)
	if err != nil {
		return company, err
	}

	return company, nil
}

func (user *Model) GetUserInfoByUserid(userId int) (GwUser, error) {
	var (
		gwuser GwUser
	)
	db, _ := utils.Connect()
	found, _ := db.Cols("name", "id", "company_id", "auth_ids").Where("id = ?", userId).Get(&gwuser)
	if found == false {
		return gwuser, fmt.Errorf("not found")
	}
	return gwuser, nil
}

func (user *Model) GetAllAuth() ([]*GwAuthority, error) {
	var (
		gwAuthList []*GwAuthority
	)
	db, _ := utils.Connect()
	if err := db.Find(&gwAuthList); err != nil {
		return gwAuthList, err
	}
	return gwAuthList, nil
}

func (user *Model) AddUser(gwuser GwUser) error {
	db, _ := utils.Connect()
	if _, err := db.Insert(&gwuser); err != nil {
		return err
	}
	return nil
}

func (user *Model) EidtUser(gwuser GwUser) (int64, error) {
	db, _ := utils.Connect()
	affected, err := db.Id(gwuser.Id).Update(gwuser)
	if err != nil {
		return affected, err
	}
	return affected, nil
}

func (user *Model) AddCom(com GwCompany) error {
	db, _ := utils.Connect()
	if _, err := db.Insert(&com); err != nil {
		return err
	}
	return nil
}

func (user *Model) EidtCom(com GwCompany) (int64, error) {
	db, _ := utils.Connect()
	affected, err := db.Id(com.Id).Update(com)
	if err != nil {
		return affected, err
	}
	return affected, nil
}
