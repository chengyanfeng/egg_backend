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
	}
	//如果不是连续登陆
	if user.LastLogInDay < (util.ToInt(the_time) - 86400) {
		user.SevenTimeSum = 1

	}

}

//喂养动作,每日奖励金币，七天连续奖励
func Feed(c *gin.Context) {
	returnp := util.P{}
	//获取userId
	userid := c.GetString("userID")
	//添加喂养动作
	feed := models.Feed{}
	feed.UserId = util.ToInt(userid)
	//添加昨天的凌晨时间
	feed.CreateTimeDay = util.GetYesDayTime()
	if !models.DB.NewRecord(feed) {

		//如果存在昨天的记录,那就获取累计天数
		models.DB.First(&feed)
		if feed.SevenTime == 7 {
			//创建feed,添加喂养记录
			feedCru := models.Feed{}
			feedCru.CreateTime = util.ToInt(time.Now().Unix())
			feed.CreateTimeDay = util.GetCurDayTime()
			feed.UserId = util.ToInt(userid)
			feed.SevenTime = 1
			models.DB.Create(&feedCru)
		} else {
			//创建feed,添加喂养记录
			feedCru := models.Feed{}
			feedCru.CreateTime = util.ToInt(time.Now().Unix())
			feed.CreateTimeDay = util.GetCurDayTime()
			feed.UserId = util.ToInt(userid)
			feed.SevenTime = feed.SevenTime + 1
			models.DB.Create(&feedCru)
			if feed.SevenTime == 7 {
				//连续登陆后添加金币
				def.AddGold(util.ToInt(userid), 1)
			}
		}
		//添加喂养记录后，添加金币
		def.AddGold(util.ToInt(userid), 1)

		//最后更改鸡的状态
		flag := def.ChangeChilckType(1, 1)
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
		feed.SevenTime = 1
		models.DB.Create(&feedCru)
		//添加金币
		def.AddGold(util.ToInt(userid), 1)

		flag := def.ChangeChilckType(1, 1)
		if flag == true {
			returnp["token"] = "token"
			returnp["code"] = def.CODE
			c.JSON(http.StatusOK, returnp)
		} else {
			returnp["token"] = "token"
			returnp["code"] = def.CODEHenErr
			c.JSON(http.StatusOK, returnp)
		}

	}

}

//用户托管
func Deposit(c *gin.Context) {

}
