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

type ResUserInfo struct {
	Id        int      `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Name      string   `json:"name" xorm:"not null default '' comment('用户名') VARCHAR(255)"`
	Phone     string   `json:"phone" xorm:"not null default '' comment('电话') unique(uniuser) VARCHAR(255)"`
	CompanyId int      `json:"company_id" xorm:"not null default 0 comment('公司ID') unique(uniuser) INT(11)"`
	Status    int      `json:"status" xorm:"not null default 1 comment('是否禁用') INT(11)"`
	Access    []string `json:"access"`
}

type GetDeviceList struct {
	Manager     string `json:"main_name"`
	Tel         string `json:"tel"`
	CompanyName string `xorm:"comname" json:"company"`
	GwDevice    `xorm:"extends"`
}

func (user *Model) GetUserList(userId int, companyId int, selectCompanyId int) ([]GetUserList, error) {
	// fmt.Println(userId)
	var (
		users []GetUserList
		// auth  []GwAuthority
	)

	if companyId != 1 && selectCompanyId != companyId {
		return users, fmt.Errorf("You have no permisson")
	}

	db, _ := utils.Connect()

	xSession := db.Table("gw_user").Select("gw_user.id,gw_user.name,gw_user.real_name,gw_user.company_id,gw_user.phone,gw_user.status,gw_user.auth_ids,gw_company.name as comname").
		Join("INNER", "gw_company", "gw_user.company_id = gw_company.id")

	if companyId == 1 && selectCompanyId == 0 {
		xSession = xSession.Where("1=1")
	} else {
		xSession = xSession.Where("gw_user.company_id = ?", selectCompanyId)
	}
	err := xSession.Find(&users)
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

func (user *Model) GetUserInfoByUserid(userId int) (ResUserInfo, error) {
	var (
		gwuser  GwUser
		auth    []GwAuthority
		resUser ResUserInfo
	)
	db, _ := utils.Connect()
	found, _ := db.Cols("name", "id", "company_id", "auth_ids").Where("id = ?", userId).Get(&gwuser)
	if found == false {
		return resUser, fmt.Errorf("not found")
	}
	db.In("id", strings.Split(gwuser.AuthIds, ",")).Find(&auth)
	for i := range auth {
		if auth[i].Access != "" {
			resUser.Access = append(resUser.Access, auth[i].Access)
		}
	}
	resUser.Id = gwuser.Id
	resUser.CompanyId = gwuser.CompanyId
	resUser.Name = gwuser.Name
	resUser.Phone = gwuser.Phone
	resUser.Status = gwuser.Status
	return resUser, nil
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
