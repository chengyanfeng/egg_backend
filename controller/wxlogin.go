package controller

import (
	"egg_backend/def"
	"egg_backend/models"
	"egg_backend/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"time"
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
	returnp := util.P{}
	defer func() {
		if err := recover(); err != nil {

			return
		}
	}()
	code, exist := c.GetQuery("code")
	if !exist {
		returnp["code"] = def.CODEWXNOCode
		c.JSON(http.StatusOK, returnp)
		return
	}
	userinfo := util.GetUserInfo(code)
	if userinfo == nil {
		return
	}
	user := models.User{}
	//用户信息是否存在
	user.WXOpenID = util.ToString((*userinfo)["openid"])
	if !models.DB.NewRecord(user) {
		//查看是否有电话号码
		models.DB.Where(user).First(&user)
		if len(user.Mobile) > 0 {

			//已经绑定过号码了，登陆成功。进行操作
			user.LoginTimes = user.LoginTimes + 1
			user.LastLogIn = util.ToInt(time.Now().Unix())
			//金币的奖励和登陆的日期
			EvenDaPrize(&user)

			returnp["code"] = def.CODE
			c.JSON(http.StatusOK, returnp)
			return
		} else {
			//没有绑定手机号，不算登陆
			returnp["code"] = def.CODENoPhone
			c.JSON(http.StatusOK, returnp)
			return
		}
	}
	//如果不存在的话，则保存用户信息
	user.WXNickName = util.ToString((*userinfo)["nickname"])
	user.WXHeadImg = util.ToString((*userinfo)["headimgurl"])
	user.CreateTime = util.ToInt(time.Now().Unix())
	models.DB.Create(&user)
	if !models.DB.NewRecord(user) {

		returnp["code"] = def.CODENoPhone
		c.JSON(http.StatusOK, returnp)
	} else {
		returnp["code"] = def.CODEErrDB
		c.JSON(http.StatusOK, returnp)
	}

	//返回用户token,从缓存里获取。
	token := util.GetCache("token")
	returnp["token"] = token
	returnp["code"] = def.CODENoPhone
	c.JSON(http.StatusOK, returnp)
	return

}

//获取验证码
func GetDentifyingCode() {
	//暂时先不写

}

//绑定手机号码
func BandPhoneNumber(c *gin.Context) {
	user := models.User{}
	returnp := util.P{}
	PhoneNumber, exist := c.GetQuery("PhoneNumber")
	token, _ := c.GetQuery("token")
	openId := util.GetCache(token)
	if !exist {
		returnp["code"] = def.CODENoPhone
		c.JSON(http.StatusOK, returnp)
	}
	//判断号码是否已经存在
	user.Mobile = PhoneNumber
	if !models.DB.NewRecord(user) {
		//查询手机的struct
		models.DB.First(&user)
		//查询wx的struct
		wxuser := models.User{}
		wxuser.WXOpenID = openId
		models.DB.First(&wxuser)
		//wx 账号绑定手机账号
		user.WXOpenID = wxuser.WXOpenID
		user.WXHeadImg = wxuser.WXHeadImg
		user.WXNickName = wxuser.WXNickName
		//保存
		models.DB.Save(&user)
		token := util.Hash("md5", user.Mobile+util.GetRandomString())
		util.AddCache(token, util.ToString(user.ID))
		returnp["token"] = token
		returnp["code"] = def.CODEPhoneBandWX
		c.JSON(http.StatusOK, returnp)
		return
	}
	//如果号码不存在，微信一定存在
	user.WXOpenID = openId
	if !models.DB.NewRecord(user) {
		models.DB.First(&user)
		user.Mobile = PhoneNumber

		//更新手机号和密码
		models.DB.Save(&user)
		//创建鸡舍
		henHouse := models.HenHouse{}
		henHouse.UserID = user.ID
		henHouse.Level = 1
		henHouse.Tools = "aaa"
		henHouse.CleanState = 1
		models.DB.Create(&henHouse)
		//创建免费公鸡
		chilck := models.Hen{}
		chilck.Name = "小么鸡"
		chilck.CreateTime = util.ToInt(time.Now().Unix())
		chilck.State = 1
		chilck.EggType = 1
		chilck.LifeTime = 365
		chilck.LifeTime = 3
		chilck.LifeValue = 3
		chilck.EggGenRate = 1
		chilck.GoldEggRate = 1
		chilck.UserID = user.ID
		chilck.HenHouseID = henHouse.ID
		models.DB.Create(&chilck)
		//把token放到缓存里面
		token := util.Hash("md5", user.Mobile+util.GetRandomString())
		util.AddCache(token, util.ToString(user.ID))
		returnp["code"] = def.CODEBandPhone
		c.JSON(http.StatusOK, returnp)
	}

}

//手机号码登陆
func PhoneNumberLogin(c *gin.Context) {

	returnp := util.P{}
	User := models.User{}
	phone := c.PostForm("phoneNumber")
	dentifyingCode := c.PostForm("dentifyingCode")
	Password, _ := c.GetQuery("PassWord")
	if len(phone) == 0 {
		returnp["code"] = def.CODEPhoneIsNull
		c.JSON(http.StatusOK, returnp)
		return
	}
	User.Mobile = phone
	if len(Password) > 0 {
		//去数据库里查询phone
		if !models.DB.NewRecord(&User) {
			models.DB.First(&User)
			if util.Hash(Password) != User.PwdHash {
				returnp["code"] = def.CODEPassWordErr
				c.JSON(http.StatusOK, returnp)
			} else {
				//返回自定义token
				util.AddCache("token", util.Hash(phone))
				returnp["code"] = def.CODE
				returnp["token"] = util.Hash(phone)
				c.JSON(http.StatusOK, returnp)
			}
		} else {
			returnp["code"] = def.CODEPhoneIsNull
			c.JSON(http.StatusOK, returnp)
		}

	} else {
		CacheDentifyCode := util.GetCache(phone)
		if len(CacheDentifyCode) == 0 {
			returnp["code"] = def.CODEDENtifyCodeExp
			c.JSON(http.StatusOK, returnp)
			return
		}
		if CacheDentifyCode != dentifyingCode {
			returnp["code"] = def.CODEDentifyCodeERR
			c.JSON(http.StatusOK, returnp)
			return
		} else {
			//如果已经注册，获取userID
			if !models.DB.NewRecord(&User) {
				models.DB.First(&User)
				//登陆次数加上一
				User.LoginTimes = User.LoginTimes + 1
				//设置最后的登陆时间
				User.LastLogIn = util.ToInt(time.Now().Unix())
				//登陆奖励金币，和七天奖励金币
				EvenDaPrize(&User)

				//把token放到缓存里面
				token := util.Hash("md5", phone+util.GetRandomString())
				util.AddCache(token, util.ToString(User.ID))
				returnp["code"] = def.CODE
				returnp["token"] = token
				c.JSON(http.StatusOK, returnp)
				return

			} else {
				//如果没有则注册手机
				User.CreateTime = util.ToInt(time.Now().Unix())
				User.LoginTimes = 1
				models.DB.Create(&User)
				//创建鸡舍
				henHouse := models.HenHouse{}
				henHouse.UserID = User.ID
				henHouse.Level = 1
				henHouse.Tools = "aaa"
				henHouse.CleanState = 1
				models.DB.Create(&henHouse)
				//创建免费公鸡
				chilck := models.Hen{}
				chilck.Name = "小么鸡"
				chilck.CreateTime = util.ToInt(time.Now().Unix())
				chilck.State = 1
				chilck.EggType = 1
				chilck.LifeTime = 365
				chilck.LifeTime = 3
				chilck.LifeValue = 3
				chilck.EggGenRate = 1
				chilck.GoldEggRate = 1
				chilck.UserID = User.ID
				chilck.HenHouseID = henHouse.ID
				models.DB.Create(&chilck)
				//把token放到缓存里面
				token := util.Hash("md5", phone+util.GetRandomString())
				util.AddCache(token, util.ToString(User.ID))
				returnp["code"] = def.CODE
				returnp["token"] = token
				c.JSON(http.StatusOK, returnp)
				return
			}
		}
	}

}

//绑定密码
func SetPassWord(c *gin.Context) {
	returnp := util.P{}
	User := models.User{}
	token, exist := c.GetQuery("token")

	Password, _ := c.GetQuery("PassWord")
	if !exist {
		returnp["code"] = def.CODETOKENErr
		c.JSON(http.StatusOK, fmt.Sprintf(" %s", returnp))
		return
	}
	userIdOrOpenId := util.GetCache(token)
	if len(userIdOrOpenId) > 12 {
		User.WXOpenID = userIdOrOpenId
		if !models.DB.NewRecord(User) {
			User.PwdHash = util.Hash(Password)
			models.DB.Save(&User)
			returnp["code"] = def.CODE
			returnp["token"] = token
			c.JSON(http.StatusOK, returnp)
		}
	} else {
		User.ID = util.ToInt(userIdOrOpenId)
		if !models.DB.NewRecord(User) {
			User.PwdHash = util.Hash(Password)
			models.DB.Save(&User)
			returnp["code"] = def.CODE
			returnp["token"] = token
			c.JSON(http.StatusOK, returnp)
		}
	}
}

//发送手机验证码
func SetNumberVcod(c *gin.Context) {
	returnp := util.P{}
	phone, exist := c.GetQuery("phoneNumber")
	if !exist {
		returnp["code"] = def.CODEPhoneIsNull
		c.JSON(http.StatusOK, returnp)
		return
	}
	ifSendOk := util.SendMessage(phone)
	if ifSendOk {
		returnp["code"] = def.CODE
		c.JSON(http.StatusOK, returnp)
		return
	} else {
		returnp["code"] = def.CODEGETERR
		c.JSON(http.StatusOK, returnp)
		return
	}

}
