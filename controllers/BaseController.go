package controllers

import (
	"gw/lib"
	"gw/model"
)

var (
	m *model.Model
)

type BaseController struct {
}

func (tc *BaseController) CheckAuth(token string, router string, ignore bool) (int, int, error) {
	claims, err := lib.ParseToken(token)
	if err != nil {
		return 0, 0, err
	}
	if ignore == true {
		return int(claims["userId"].(float64)), int(claims["comId"].(float64)), nil
	}

	if authErr := m.CheckAuthByUserId(claims["userId"].(float64), router); authErr != nil {
		return 0, 0, authErr
	}
	return int(claims["userId"].(float64)), int(claims["comId"].(float64)), nil
}
