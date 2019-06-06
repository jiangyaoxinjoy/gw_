package utils

import (
	"fmt"
	"gw/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func Connect() (*xorm.Engine, error) {
	db, err := xorm.NewEngine("mysql", config.Mysql)
	if err != nil {
		return nil, fmt.Errorf("Mysql Error:" + err.Error())
	}
	return db, err
}
