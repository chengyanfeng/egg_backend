package controller

import (
	"egg_backend/def"
	"egg_backend/models"
	"egg_backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//每日登录奖励，每日奖励金币，七天连续奖励
func EvenDaPrize(user *models.User) {
	the_time, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	//登陆后金币添加+1
	userProper := models.UserProperty{}
	userProper.UserID = user.ID
	henHouse := models.HenHouse{}
	henHouse.UserID = user.ID
	//添加金币
	if !models.DB.NewRecord(userProper) {
		//查询UserProperty
		models.DB.First(&userProper)
		//金币添加1
		userProper.Coins = userProper.Coins + 1
		models.DB.Save(&userProper)
	}
	//判断是否连续登陆
	if user.LastLogInDay == (util.ToInt(the_time) - 86400) {
		//如果前一天登陆过，则把现在日期赋值给当前日期，七天累加登陆+1
		user.SevenTimeSum = user.SevenTimeSum + 1
		user.LastLogInDay = util.ToInt(the_time)
		//判断七天的连续登陆
		if user.SevenTimeSum == 7 {

			//添加金币，然后回归为user表的SevenTimeSum 七天累计天数为1
			user.SevenTimeSum = 1

			//添加金币
			if !models.DB.NewRecord(userProper) {
				//查询UserProperty
				models.DB.First(&userProper)
				//金币添加1
				userProper.Coins = userProper.Coins + 1
				models.DB.Save(&userProper)
			}

		}
		//判断三天 的连续登陆
		if user.SevenTimeSum == 3 {
			//查询鸡的鸡舍
			models.DB.First(&henHouse)
			//查询鸡舍是否升级到二级
			if henHouse.Level < 2 {
				henHouse.Level = 2
				//更新鸡舍等级
				models.DB.Model(&henHouse).Update("level")
				//获得一枚金蛋
				egg := models.Egg{}
				egg.Type = 5
				egg.Type = util.ToInt(time.Now().Unix())
				egg.UserId = user.ID
				egg.CreateTimeDay = util.GetCurDayTime()
				egg.Sell = 0
				egg.Source = 3
				//添加金蛋记录
				models.DB.Create(&egg)
				//个人资产添加记录
				userProperty := models.UserProperty{}
				userProperty.UserID = user.ID
				models.DB.First(&userProperty)
				userProperty.GoldEggsFree = userProperty.GoldEggsFree + 1
				//更新资产记录
				models.DB.Model(&userProperty).Update("GoldEggsFree")
			}
		}
	}
	//如果不是连续登陆
	if user.LastLogInDay < (util.ToInt(the_time) - 86400) {
		user.SevenTimeSum = 1

	}

}

//喂养动作,每日奖励金币，七天连续奖励
func Feed(c *gin.Context) {
	returnp := util.P{}
	token := c.GetString("token")
	//获取userId
	userid := util.GetCache(token)
	HenId := c.GetString("henId")
	feedType := c.GetString("feedType")
	//添加喂养动作
	feed := models.Feed{}
	feed.UserId = util.ToInt(userid)
	//添加昨天的凌晨时间
	feed.CreateTimeDay = util.GetYesDayTime()
	if !models.DB.NewRecord(feed) {
		//判断七天是否连续喂养，然后处理
		flag := def.SevenRecord(feed, userid, HenId, util.ToInt(feedType))
		if flag == true {
			returnp["token"] = "token"
			returnp["code"] = def.CODE
			c.JSON(http.StatusOK, returnp)
		} else {
			returnp["token"] = "token"
			returnp["code"] = def.CODEHenErr
			c.JSON(http.StatusOK, returnp)
		}
	} else {
		//创建feed,添加喂养记录
		feedCru := models.Feed{}
		feedCru.CreateTime = util.ToInt(time.Now().Unix())
		feed.CreateTimeDay = util.GetCurDayTime()
		feed.UserId = util.ToInt(userid)
		feed.HenId = util.ToInt(HenId)
		feed.SevenTime = 1
		models.DB.Create(&feedCru)
		//添加金币
		def.AddGold(util.ToInt(userid), 1)

		flag := def.ChangeChilckType(1, 1)
		if flag == true {
			returnp["token"] = token
			returnp["code"] = def.CODE
			c.JSON(http.StatusOK, returnp)
		} else {
			returnp["token"] = token
			returnp["code"] = def.CODEHenErr
			c.JSON(http.StatusOK, returnp)
		}

	}

}

//用户托管
func Deposit(c *gin.Context) {
	returnp := util.P{}
	token := c.GetString("token")
	henId := c.GetString("henId")
	//获取userId
	userid := util.GetCache(token)
	hen := models.Hen{}
	hen.ID = util.ToInt(henId)
	hen.UserID = util.ToInt(userid)
	if !models.DB.NewRecord(hen) {
		models.DB.First(&hen)
		hen.State = 2
		models.DB.Model(&hen).Update("state")
		returnp["token"] = token
		returnp["code"] = def.CODE
		c.JSON(http.StatusOK, returnp)
	} else {
		returnp["token"] = token
		returnp["code"] = def.CODEDepositErr
		c.JSON(http.StatusOK, returnp)
	}

}

//鸡状态的获取

func GetHenNature(c *gin.Context) {
	returnp := util.P{}
	token := c.GetString("token")
	henId := c.GetString("henId")
	//获取userId
	userid := util.GetCache(token)
	hen := models.Hen{}
	hen.ID = util.ToInt(henId)
	hen.UserID = util.ToInt(userid)
	if !models.DB.NewRecord(hen) {
		models.DB.First(&hen)
		returnp["token"] = token
		returnp["code"] = def.CODE
		returnp["henId"] = hen
		c.JSON(http.StatusOK, returnp)
	} else {
		returnp["token"] = token
		returnp["code"] = def.CODEHenErr
		c.JSON(http.StatusOK, returnp)
	}

}
