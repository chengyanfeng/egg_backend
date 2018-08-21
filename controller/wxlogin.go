package controller

import (
	"github.com/gin-gonic/gin"
	"Egg/def"
	"net/http"
	"fmt"
	"Egg/util"
)

//跳转微信的url
func RedirecturlHandler(c *gin.Context) {
	redirecturl := "https%3a%2f%2fservice.genyuanlian.com%2fseven_night%2findex"
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + def.WEIXINAPPID + "&redirect_uri=" + redirecturl +
		"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"
	c.Redirect(302, url)
}

//微信跳转回来的url
func IndexHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Redirect(302, "/wxlogin/url", )
		}
	}()
	code, exist := c.GetQuery("code")
	if !exist {
		code = "the code is not exist!"
		c.JSON(http.StatusOK, []byte(fmt.Sprintf("%s", code)))
		return
	}
	userinfo := util.GetUserInfo(code)
	if userinfo == nil {
		c.Redirect(302, "/wxlogin/url")
	}
	//返回用户token
	token:=util.GetCache("token")
	c.JSON(http.StatusOK,  []byte(fmt.Sprintf("%s", token)))
}

//手机号码登陆

