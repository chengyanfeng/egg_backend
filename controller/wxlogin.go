package controller

import (
	"egg_backend/def"
	"egg_backend/models"
	"egg_backend/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//跳转微信的url,暂时无用，改为前端调用
func RedirectUrlHandler(c *gin.Context) {
	redirecturl := "https%3a%2f%2fservice.genyuanlian.com%2fseven_night%2findex"
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + def.WEIXINAPPID + "&redirect_uri=" + redirecturl +
		"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"
	c.Redirect(302, url)
}

//微信跳转回来的url
func IndexHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Redirect(302, "/wxlogin/url")
			return
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
		return
	}
	//返回用户token
	token := util.GetCache("token")
	c.JSON(http.StatusOK, []byte(fmt.Sprintf("%s", token)))
	return

}

//获取验证码
// 14c938c682dfadee623ab52dbda8d5e3f086617d
func GetDentifyingCode() {
	//暂时先不写

}

//绑定手机号码
func BandPhoneNumber(c *gin.Context) {
	value, exist := c.GetQuery("PhoneNumber")
	token, _ := c.GetQuery("token")
	if !exist {
		value = "the PhoneNumber is not exist!"
		c.JSON(http.StatusOK, []byte(fmt.Sprintf(" %s", value)))
	}
	//获取openId
	openId := util.GetCache(token)
	fmt.Print(openId)
	//根据openId像表中插入手机号码

}

//手机号码登陆
func PhoneNumberLogin(c *gin.Context) {
	User := models.User{}
	hose := models.HenHouse{}
	hose.UserID = 890
	User.ID = 32424
	User.Mobile = "1232132132132"
	//去数据库里查询phone
	models.DB.Create(&User)
	phone, exist := c.GetQuery("phoneNumber")
	dentifyingCode, exitstCode := c.GetQuery("dentifyingCode")
	if !exist {
		phone = "the PhoneNumber is not exist!"
		c.JSON(http.StatusOK, fmt.Sprintf(" %s", phone))
		return
	}
	if !exitstCode {
		dentifyingCode = "the dentifyingCode is not exist!"
		c.JSON(http.StatusOK, []byte(fmt.Sprintf(" %s", dentifyingCode)))
		return
	}
	CacheDentifyCode := util.GetCache(phone)
	if CacheDentifyCode != dentifyingCode {
		value := "the dentifyingCode is not match!"
		c.JSON(http.StatusOK, []byte(fmt.Sprintf(" %s", value)))
		return
	}

	//返回自定义token

}
