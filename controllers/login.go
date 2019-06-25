package controllers

import (
	"gw/lib"
	"gw/model"
	"gw/utils"

	"time"

	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// LoginResult 登录结果结构
type LoginResult struct {
	Token string       `json:"token"`
	User  model.GwUser `json:"user"`
}

//登录
func (tc *BaseController) Login(c *gin.Context) {
	var loginReq model.LoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "json 解析失败",
		})
		return
	}
	loginReq.Password = utils.String2md5(loginReq.Password)
	user, err := m.LoginCheck(loginReq)
	if err != nil {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "验证失败:" + err.Error(),
		})
		return
	}
	fmt.Println(user.Id)
	if token, err := generateToken(user); err == nil {
		data := LoginResult{
			Token: token,
			User:  user,
		}
		c.JSON(200, gin.H{
			"status": 0,
			"msg":    "登录成功！",
			"data":   data,
		})
		return
	} else {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "验证失败," + err.Error(),
		})
		return
	}

}

func (tc *BaseController) Home(c *gin.Context) {
	fmt.Println(c.Request.RequestURI)
	token := c.Query("token")
	token = strings.Trim(token, " ")
	if token == "" {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "没有token",
		})
		return
	}
	claims, err := lib.ParseToken(token)
	if err != nil {
		if token == "" {
			c.JSON(200, gin.H{
				"status": -1,
				"msg":    "token错误",
			})
			return
		}
	}
	expires := claims["expires"].(float64)

	if int64(expires) >= int64(time.Now().Unix()) {
		c.JSON(200, gin.H{
			"status": 0,
			"msg":    "token正常",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"status": -1,
			"msg":    "token过期",
		})
		return
	}
}

func (tc *BaseController) ChangePsd(c *gin.Context) {
	var (
		token model.ReqChangePsd
	)
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	userId, comId, autherr := tc.CheckAuth(token.TokenString, "", true)
	if autherr != nil {
		c.JSON(200, gin.H{"status": -1, "msg": autherr.Error()})
		return
	}
	if token.NewPsd != token.RepeatPsd {
		c.JSON(200, gin.H{"status": -1, "msg": fmt.Errorf("输入的密码不一致")})
		return
	}
	if userId != token.UserId {
		c.JSON(200, gin.H{"status": -1, "msg": fmt.Errorf("用户错误")})
		return
	}
	err := m.ChangePassword(token, comId)
	if err != nil {
		c.JSON(200, gin.H{"status": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": 0, "msg": "OK", "data": ""})
	return

}

// 生成令牌
func generateToken(user model.GwUser) (string, error) {
	claims := jwt.MapClaims{
		"userId":  int(user.Id),
		"comId":   user.CompanyId,
		"expires": int64(time.Now().Unix() + 3600*24), // 过期时间 24小时
	}
	token, err := lib.CreateToken(claims)
	if err != nil {
		return token, err
	}
	return token, nil
}
