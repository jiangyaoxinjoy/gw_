package controllers

import (
	"gw/model"

	"github.com/gin-gonic/gin"
)

type ReqAllAlert struct {
	ReqMonitoring
	// State string `json:"state" binding:"required"`
}

type ReqStateAlert struct {
	ReqAllAlert
	State string `json:"state" binding:"required"`
}

func (tc *BaseController) WxAlertAllCount(c *gin.Context) {
	var (
		query ReqAllAlert
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	// token := c.Request.Header.Get("token")
	_, companyId, authErr := tc.CheckAuth(query.TokenString, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	data, err := m.WxAllAlertCount(query.CompanyId, companyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) StateAlert(c *gin.Context) {
	var (
		query ReqStateAlert
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	// token := c.Request.Header.Get("token")
	_, companyId, authErr := tc.CheckAuth(query.TokenString, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	data, err := m.WxStateAlert(query.CompanyId, companyId, query.State)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}

func (tc *BaseController) WxSetupDevice(c *gin.Context) {
	var (
		query model.SetupDevice
	)
	router := c.Request.RequestURI
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}

	_, _, authErr := tc.CheckAuth(query.TokenString, router, false)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	if err := m.WxSetupDevice(query); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK"})
	return
}
