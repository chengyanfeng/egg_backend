package controller

import (
	"egg_backend/models"
	"egg_backend/util"
)

//每个小鸡每天的生命状态,定时跑批的，如果喂养就下蛋，如果没有喂养就减去生命值
func HenLiveCheck() {
	//获取所有的小鸡列表
	HenList := []models.Hen{}
	models.DB.Where("life_time > ?", 0).Find(&HenList)
	for _, v := range HenList {
		//所有的鸡全部更为饥饿状态
		v.State = 1
		models.DB.Model(&v).Update("state")

		feed := models.Feed{}
		feed.HenId = v.ID
		//检查昨天是否喂养
		feed.CreateTimeDay = util.GetYesDayTime()
		//如果喂养,下蛋
		if !models.DB.NewRecord(feed) {
			//鸡舍
			HenHouse := models.HenHouse{}
			HenHouse.UserID = v.UserID
			models.DB.First(&HenHouse)
			//用户资产表
			userProperty := models.UserProperty{}
			userProperty.UserID = v.UserID
			models.DB.First(&userProperty)
			//查询喂养记录
			models.DB.First(&feed)
			//判断鸡的类型
			switch v.HenType {
			case 1:
				models.FreeHee(&feed, &userProperty, &v, &HenHouse)
				break
			case 2:
				models.DarlingHee(&feed, &userProperty, &v, &HenHouse)
				break
			case 3:
				models.GoldHee(&feed, &userProperty, &v, &HenHouse)
				break

			case 4:
				models.LitileHee(&feed, &userProperty, &v, &HenHouse)
				break
			}

		} else {
			//没有喂养公鸡生命值-1
			v.Lifes = v.Lifes - 1
			//公鸡生命值减去1
			models.DB.Save(&v)
		}
		//查看今天的是否喂养
		feedif := models.Feed{}
		feed.HenId = v.ID
		feed.UserId = v.UserID
		if !models.DB.NewRecord(feedif) {
			//更改鸡的状态
			models.ChangeChilckType(v.ID, 2)
		}

		//如果是托管的鸡
		if v.Deposit == 1 {
			//查看当天有没有喂养
			goldFeed := models.Feed{}
			count := 0
			models.DB.Model(&goldFeed).Where("UserId=? AND HenId=? AND CreateTimeDay=?", v.UserID, v.ID, util.ToInt(util.GetCurDayTime())).Count(&count)
			if count == 0 {
				//主动喂食料
				UserProperty := models.UserProperty{}
				UserProperty.UserID = v.UserID
				models.DB.First(&UserProperty)
				//检查是否有还有食料
				flag := models.CheckFeed(&UserProperty)
				if flag > 0 {
					models.DepositFeed(flag, &v)
				}
			}
		}

	}

}
