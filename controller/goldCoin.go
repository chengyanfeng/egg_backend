package controller

import (
	"egg_backend/models"
	"egg_backend/util"
	"time"
)

//每日登录奖励
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
