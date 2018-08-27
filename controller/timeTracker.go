package controller

import (
	"egg_backend/def"
	"egg_backend/models"
	"egg_backend/util"
)

//每个小鸡每天的生命状态,定时跑批的，如果喂养就下蛋，如果没有喂养就减去生命值
func HenLiveCheck() {
	//获取所有的小鸡列表
	HenList := []models.Hen{}
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
			//已经喂养，生产鸡蛋。
			models.DB.First(&feed)
			//判断鸡的类型
			switch v.HenType {
			case 1:
				def.FreeHee(&feed, &userProperty, &v, &HenHouse)
				break
			case 2:
				def.DarlingHee(&feed, &userProperty, &v, &HenHouse)
				break
			case 3:
				def.DarlingHee(&feed, &userProperty, &v, &HenHouse)
				break
			}

		} else {
			//没有喂养公鸡生命值-1
			v.Lifes = v.Lifes - 1
			//公鸡生命值减去1
			models.DB.Save(&v)
		}
		//如果是托管的公鸡
		if v.Deposit == 1 {
			//主动喂食料

			UserProperty := models.UserProperty{}
			UserProperty.UserID = v.UserID
			models.DB.First(&UserProperty)
			//检查是否有还有食料
			flag := def.CheckFeed(&UserProperty)
			if flag > 0 {
				def.DepositFeed(flag, v)
			}
		}

	}

}
