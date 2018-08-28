package controller

import (
	"egg_backend/def"
	"egg_backend/models"
	"egg_backend/util"
	"fmt"
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
		models.FeedType(util.ToInt(feedType), util.ToInt(userid), util.ToInt(HenId))
		flag := models.ChangeChilckType(util.ToInt(HenId), 2)
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
		models.FeedType(util.ToInt(feedType), util.ToInt(userid), util.ToInt(HenId))
		flag := models.ChangeChilckType(util.ToInt(HenId), 2)
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

//产蛋数量的查询
func GetEggNumber(c *gin.Context) {
	a := 0
	returnp := util.P{}
	token := c.GetString("token")
	//获取userId
	userid := util.GetCache(token)
	henId := c.GetString("henId")
	fmt.Print(henId)
	egg := models.Egg{}
	models.DB.Model(&egg).Where("UserId=? AND HenId=?", userid, henId).Count(&a)
	returnp["token"] = token
	returnp["code"] = def.CODE
	returnp["count"] = a
	c.JSON(http.StatusOK, returnp)
}

//系统兑换鸡蛋
func EggShell(c *gin.Context) {
	returnp := util.P{}
	token := c.GetString("token")
	//获取userId
	userid := util.GetCache(token)
	redEgg := c.GetString("redEgg")
	realEgg := c.GetString("realEgg")

	if len(redEgg) > 0 {
		//查询是否有那么多彩蛋
		userProperty := models.UserProperty{}
		userProperty.UserID = util.ToInt(userid)
		models.DB.First(&userProperty)
		if userProperty.FreeEggs < util.ToInt(redEgg) {
			returnp["token"] = token
			returnp["code"] = def.CODENumberErr
			c.JSON(http.StatusOK, returnp)
		}
		if util.ToInt(redEgg) > 0 {
			egg := []models.Egg{}
			//查询彩蛋
			models.DB.Where("UserId=? AND Type=? AND Sell=?", userid, 1, 0).Find(&egg)
			//根据彩蛋兑换的数量更新彩蛋的售卖状态
			for k, v := range egg {
				if k%3 == 2 {
					//如果满三次，那就换成鸡蛋
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//换成鸡蛋
					egg := models.Egg{}
					egg.Sell = 0
					egg.CreateTimeDay = util.GetCurDayTime()
					egg.CreateTime = util.ToInt(time.Now().Unix())
					egg.UserId = util.ToInt(userid)
					egg.Type = 2
					egg.Source = 2

					//个人资产彩蛋减少一个,鸡蛋添加一个
					userProperty.FreeEggs = userProperty.FreeEggs - 1
					userProperty.RealEggs = userProperty.RealEggs + 1
					models.DB.Model(&userProperty).Update("FreeEggs", "RealEggs")

					if k == util.ToInt(redEgg)-1 {
						break
					}

				} else {
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//个人资产彩蛋减少一个
					userProperty.FreeEggs = userProperty.FreeEggs - 1
					models.DB.Model(&userProperty).Update("FreeEggs")
				}
			}
		}
	}
	if len(realEgg) > 0 {
		//查询是否有那么多鸡蛋
		userProperty := models.UserProperty{}
		userProperty.UserID = util.ToInt(userid)
		models.DB.First(&userProperty)
		if userProperty.RealEggs < util.ToInt(realEgg) {
			returnp["token"] = token
			returnp["code"] = def.CODENumberErr
			c.JSON(http.StatusOK, returnp)
		}
		if util.ToInt(realEgg) > 0 {
			egg := []models.Egg{}
			//查询鸡蛋
			models.DB.Where("UserId=? AND Type=? AND Sell=?", userid, 2, 0).Find(&egg)
			//根据鸡蛋回收
			for k, v := range egg {
				if k%10 == 9 {
					//如果满10次，那就兑换成金币
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//系统回收表添加记录
					EggWithdrawOrder := models.EggWithdrawOrder{}
					EggWithdrawOrder.CreateTime = util.ToInt(time.Now().Unix())
					EggWithdrawOrder.UserID = util.ToInt(userid)
					EggWithdrawOrder.Amount = 10
					EggWithdrawOrder.PriceCent = 5
					models.DB.Create(&EggWithdrawOrder)
					//个人资产鸡蛋减少一个,金币添加5个
					userProperty.RealEggs = userProperty.RealEggs - 1
					userProperty.Coins = userProperty.Coins + 5
					models.DB.Model(&userProperty).Update("RealEggs", "Coins")

					if k == util.ToInt(redEgg)-1 {
						break
					}

				} else {
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//个人资产鸡蛋减少一个
					userProperty.RealEggs = userProperty.RealEggs - 1
					models.DB.Model(&userProperty).Update("RealEggs")
				}
			}
		}
	}

	returnp["token"] = token
	returnp["code"] = def.CODE
	c.JSON(http.StatusOK, returnp)
}

//系统提现鸡蛋
func EggExchange(c *gin.Context) {
	returnp := util.P{}
	token := c.GetString("token")
	//获取userId
	userid := util.GetCache(token)
	realEgg := c.GetString("realEgg")

	if len(realEgg) > 0 {
		//查询是否有那么鸡蛋
		userProperty := models.UserProperty{}
		userProperty.UserID = util.ToInt(userid)
		models.DB.First(&userProperty)
		if userProperty.RealEggs < util.ToInt(realEgg) {
			returnp["token"] = token
			returnp["code"] = def.CODENumberErr
			c.JSON(http.StatusOK, returnp)
		}
		if util.ToInt(realEgg) > 0 {
			egg := []models.Egg{}
			//查询鸡蛋
			models.DB.Where("UserId=? AND Type=? AND Sell=?", userid, 2, 0).Find(&egg)
			//根据鸡蛋兑换的
			for k, v := range egg {
				if k%30 == 29 {
					//如果满30个
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//个人资产鸡蛋减少一个
					userProperty.RealEggs = userProperty.RealEggs - 1
					models.DB.Model(&userProperty).Update("RealEggs")
					//添加邮寄信息
					eggTakenOrder := models.EggTakenOrder{}
					eggTakenOrder.UserID = util.ToInt(userid)
					eggTakenOrder.State = 1
					eggTakenOrder.CreateTime = util.ToInt(time.Now().Unix())
					eggTakenOrder.Amount = 30
					eggTakenOrder.DeliverInfo = "xxx"
					user := models.User{}
					user.ID = util.ToInt(userid)
					models.DB.First(&user)
					eggTakenOrder.Address = user.AddressInfo
					models.DB.Create(&eggTakenOrder)
					if k == util.ToInt(realEgg)-1 {
						break
					}

				} else {
					v.Sell = 1
					//更新状态，标记已经售卖，或者兑换
					models.DB.Model(&v).Update("Sell")
					//个人资产鸡蛋减少一个
					userProperty.RealEggs = userProperty.RealEggs - 1
					models.DB.Model(&userProperty).Update("RealEggs")
				}
			}
		}
	}

	returnp["token"] = token
	returnp["code"] = def.CODE
	c.JSON(http.StatusOK, returnp)
}
