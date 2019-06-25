package model

import (
	"time"
)

type GwAlert struct {
	Id        int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	DeviceId  string `json:"device_id" xorm:"comment('设备ID') VARCHAR(255)"`
	MessageId string `json:"message_id" xorm:"not null default '' comment('报警编号') VARCHAR(255)"`
	AlertType string `json:"alert_type" xorm:"not null default '10' comment('报警类型:10=压力,20=偷水,30=撞到,40=在线,50=信号强度') ENUM('10','20','30','40','50')"`
	Cola      string `json:"cola" xorm:"comment('通知参数1') VARCHAR(255)"`
	Colb      string `json:"colb" xorm:"comment('通知参数2') VARCHAR(255)"`
	Colc      string `json:"colc" xorm:"comment('通知参数3') VARCHAR(255)"`
	Totala    string `json:"totala" xorm:"comment('在线') VARCHAR(255)"`
	Totalb    string `json:"totalb" xorm:"comment('偷水') VARCHAR(255)"`
	Totalc    string `json:"totalc" xorm:"comment('撞倒') VARCHAR(255)"`
	Totald    string `json:"totald" xorm:"comment('开机') VARCHAR(255)"`
	// Totale    string `json:"totale" xorm:"comment('水压') VARCHAR(255)"`
	// Totalf     string    `json:"totalf" xorm:"comment('信号') VARCHAR(255)"`
	Pstate      int       `json:"pstate" xorm:"not null default 0 comment(''标记水压异常：1=异常，0=恢复)" INT(11)"`
	Createtime  time.Time `json:"createtime" xorm:"-"`
	CompanyId   int       `json:"company_id" xorm:"not null default 0 comment('公司ID') INT(11)"`
	Descrip     string    `json:"descrip" xorm:"comment('备注') VARCHAR(255)"`
	Sendtime    string    `json:"sendtime" xorm:"comment('发送时间') VARCHAR(255)"`
	Restoretime int       `json:"restoretime" xorm:"comment('水压恢复时间') INT(11)"`
	RestoreId   int       `json:"restore_id" xorm:"INT(11)"`
}

type GwAuthSub struct {
	Id     int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Node   string `json:"node" xorm:"VARCHAR(255)"`
	AuthId int    `json:"auth_id" xorm:"not null default 0 INT(11)"`
	Name   string `json:"name" xorm:"VARCHAR(255)"`
}

type GwAuthority struct {
	Id     int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Name   string `json:"name" xorm:"not null default '' comment('权限名称') VARCHAR(255)"`
	Access string `json:"access" xorm: "VARCHAR(255)"`
}

type GwCompany struct {
	Id         int       `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Name       string    `json:"name" xorm:"not null default '' comment('公司名称') unique VARCHAR(255)"`
	Address    string    `json:"address" xorm:"default '' comment('公司地址') VARCHAR(255)"`
	Value1     string    `json:"value1" xorm:"not null default '0.2' comment('压力阀值1') VARCHAR(255)"`
	Value2     string    `json:"value2" xorm:"not null default '0.35' comment('压力阀值2') VARCHAR(255)"`
	Createtime time.Time `json:"createtime" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Email      string    `json:"email" xorm:"not null default '' VARCHAR(255)"`
	Tel        string    `json:"tel" xorm:"not null default '' VARCHAR(255)"`
	Manager    string    `json:"manager" xorm:"not null default '' VARCHAR(255)"`
}

type GwDevice struct {
	Id         int       `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Address    string    `json:"address" xorm:"comment('设备地址') VARCHAR(255)"`
	Lng        string    `json:"lng" xorm:"comment('经度') VARCHAR(255)"`
	Lat        string    `json:"lat" xorm:"comment('纬度') VARCHAR(255)"`
	DeviceId   string    `json:"device_id" xorm:"not null default '' comment('设备号') unique VARCHAR(255)"`
	State      string    `json:"state" xorm:"comment('当前设备状态') VARCHAR(255)"`
	CompanyId  int       `json:"company_id" xorm:"not null default 0 comment('所属公司ID') INT(11)"`
	Status     int       `json:"status" xorm:"comment('设备是否安装') INT(11)"`
	Createtime time.Time `json:"createtime" xorm:"-"`
	Setuptime  string    `json:"setuptime" xorm:"comment('安装时间') VARCHAR(255)"`
	Hearttime  string    `json:"hearttime" xorm:"comment('心跳时间') VARCHAR(255)"`
	AlertId    int       `json:"alert_id" xorm:"INT(11)"`
	Signal     string    `json:"signal" xorm:"VARCHAR(255)"`
	Descrip    string    `json:"descrip" xorm:"VARCHAR(255)"`
}

type GwPressure struct {
	Id            int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	CompanyId     int    `json:"company_id" xorm:"not null default 0 comment('公司ID') INT(11)"`
	DeviceId      string `json:"device_id" xorm:"comment('设备ID') VARCHAR(255)"`
	Sendtime      string `json:"sendtime" xorm:"comment('发送时间') VARCHAR(255)"`
	PressureValue string `json:"pressure_value" xorm:"comment('压力值') VARCHAR(255)"`
	MsgId         int    `json:"msg_id" xorm:"INT(11)" binding:"required"`
}

type GwUser struct {
	Id         int       `json:"Id" xorm:"not null pk autoincr INT(11)"`
	RealName   string    `json:"real_name" xorm:"not null default '' comment('真实姓名') VARCHAR(255)"`
	Name       string    `json:"name" xorm:"not null default '' comment('用户名') VARCHAR(255)"`
	Password   string    `json:"password" xorm:"not null default '' comment('密码') VARCHAR(255)"`
	Phone      string    `json:"phone" xorm:"not null default '' comment('电话') unique(uniuser) VARCHAR(255)"`
	CompanyId  int       `json:"company_id" xorm:"not null default 0 comment('公司ID') unique(uniuser) INT(11)"`
	AuthIds    string    `json:"auth_ids" xorm:"default '' comment('权限ID:1=增加,2=删除,3=更新设备信息,4=安装设备,5=接收通知') VARCHAR(255)"`
	Createtime time.Time `json:"createtime" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	LoginTime  time.Time `json:"login_time" xorm:"comment('登录时间') TIMESTAMP"`
	Status     int       `json:"status" xorm:"not null default 1 comment('是否禁用') INT(11)"`
}
type GwNotify struct {
	Id          int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	UserId      int    `json:"user_id" xorm:"comment('通知人ID') INT(11)"`
	AlertId     int    `json:"alert_id" xorm:"comment('报警ID') INT(11)"`
	Type        int    `json:"type" xorm:"comment('通知方式:1=短信,2=微信') INT(11)"`
	Sendtime    int    `json:"sendtime" xorm:"comment('通知发送时间') INT(11)"`
	DeviceId    string `json:"device_id" xorm:"comment('设备ID') VARCHAR(255)"`
	State       int    `json:"state" xorm:"INT(11)"`
	Receivetime int    `json:"receivetime" xorm:"INT(11)"`
}
