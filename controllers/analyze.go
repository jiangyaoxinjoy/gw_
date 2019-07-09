package controllers

import (
	"gw/model"

	"github.com/gin-gonic/gin"
)

func (tc *BaseController) AlertAnalyze(c *gin.Context) {
	var (
		params model.ReqAnalyze
	)
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	token := c.Request.Header.Get("token")
	_, companyId, authErr := tc.CheckAuth(token, "", true)
	if authErr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": authErr.Error()})
		return
	}
	data, err := m.GetAlertAnalyze(params, companyId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": data})
	return
}
