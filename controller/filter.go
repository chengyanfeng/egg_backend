package controller

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"egg_backend/util"
	"net/http"
)

//登陆拦截器
func LoginFilter(c *gin.Context) {
	if c.Request.URL.Path == "/wxlogin/url" || c.Request.URL.Path == "/wxlogin/index" ||c.Request.URL.Path=="/PhoneNumberLogin"{
		fmt.Print(c.Request.URL.Path)
		c.Next()
	} else {
		value, exist := c.GetQuery("token")
		if !exist {
			value = "the token is not exist!"
			c.JSON(http.StatusOK, []byte(fmt.Sprintf(" %s", value)))
			return
		}
		openId := util.GetCache(value)
		if openId == "" {
			c.JSON(http.StatusOK, []byte(fmt.Sprintf("%s", "token is  expire")))
			return
		} else {
			c.Next()
		}

	}

}
