package def

import (
	"egg_backend/models"
	"egg_backend/util"
	"math"
	"math/rand"
	"time"
)

//添加金币
func AddGold(userId, number int) {
	propertyEvenDay := models.UserProperty{}
	propertyEvenDay.UserID = userId
	models.DB.First(&propertyEvenDay)
	propertyEvenDay.Coins = propertyEvenDay.Coins + number
	models.DB.Model(&propertyEvenDay).Update("coins")
}

//更改鸡的状态
func ChangeChilckType(HenId, Type int) bool {
	if !models.DB.NewRecord(HenId) {
		Hen := models.Hen{}
		Hen.ID = HenId
		models.DB.First(&Hen)
		Hen.State = Type
		models.DB.Model(&Hen).Update("state")
		return true
	} else {
		//返回错误
		return false
	}
}

//鸡的生蛋随机概率，根据鸡的鸡舍+根据鸡食料，概率累加，超整数，下蛋
func EggFeedHouseProbability(level int, feedNumber int) (Probability float64) {
	FeedFloat := util.ToFloat(feedNumber)
	HouseProbability1 := []float64{0.6, 0.7, 0.8, 0.9, 1}
	HouseProbability2 := []float64{0.6, 0.75, 0.8, 0.9, 1}
	HouseProbability3 := []float64{0.65, 0.75, 0.85, 0.95, 1.05}
	HouseProbability4 := []float64{0.7, 0.8, 0.9, 1, 1.1}
	HouseProbability5 := []float64{0.7, 0.8, 0.9, 1, 1.1}
	t := rand.Intn(4)

	switch level {
	case 1:
		Probability = HouseProbability1[t] * FeedFloat
		break
	case 2:
		Probability = HouseProbability2[t] * FeedFloat
		break
	case 3:
		Probability = HouseProbability3[t] * FeedFloat
		break
	case 4:
		Probability = HouseProbability4[t] * FeedFloat
		break
	case 5:
		Probability = HouseProbability5[t] * FeedFloat
		break
	default:
		Probability = HouseProbability1[t] * FeedFloat
		break

	}
	return
}

//鸡生金蛋的概率，根据鸡舍，
func EggGold(v *models.Hen, henHouse *models.HenHouse) (randint int) {
	randint = 0
	//乖乖鸡
	if v.HenType == 2 {
		//鸡舍等级，判断金蛋概率
		switch henHouse.Level {
		case 2:
			//百分之1到百分之五
			randint = rand.Intn(4000) + 1000
			break
			//百分之5到百分之10
		case 3:
			randint = rand.Intn(500) + 500
			break
			//百分之5到百分之10
		case 4:
			randint = rand.Intn(500) + 500
			break
			//百分之10-15
		case 5:
			randint = rand.Intn(177) + 333
			break
		}
	}
	//金鸡
	if v.HenType == 2 {
		//鸡舍等级，判断金蛋概率
		switch henHouse.Level {
		case 2:
			//百分之5到百分之10
			randint = rand.Intn(500) + 500
			break
			//百分之10到百分十15
		case 3:
			randint = rand.Intn(177) + 333
			break
			//百分之10到百分十15
		case 4:
			randint = rand.Intn(177) + 333
			break
			//百分之15到20
		case 5:
			randint = rand.Intn(133) + 200
			break
		}
	}
	return randint
}

//获取免费鸡生蛋跟据喂养的饲料，根据概率，产生的蛋
func FreeHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, henHouse *models.HenHouse) {
	//如果食料为系统赠送的饲料
	if feed.Type == 1 {
		heeAgg(feed, userProperty, v, 1, henHouse)
	}

	//如果是普通的饲料
	if feed.Type == 2 {
		rate := math.Floor(v.EggGenRate)
		v.EggGenRate = v.EggGenRate + EggFeedHouseProbability(henHouse.Level, feed.Type)
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			heeAgg(feed, userProperty, v, number, henHouse)
		}
	}

}

//获取乖乖鸡生蛋跟据喂养的饲料，根据概率，产生的蛋
func DarlingHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, henHouse *models.HenHouse) {

	rate := math.Floor(v.EggGenRate)
	v.EggGenRate = v.EggGenRate + EggFeedHouseProbability(henHouse.Level, feed.Type)
	nowRate := math.Floor(v.EggGenRate)
	number := util.ToInt(nowRate - rate)
	if rate == nowRate {
		//概率没有超过整数不下蛋
	} else {
		//根据概率获取下鸡蛋的个数
		heeAgg(feed, userProperty, v, number, henHouse)
	}

}

//获取,金鸡生蛋跟据喂养的饲料，根据概率，根据鸡舍，产生的蛋
func GoldHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, henHouse *models.HenHouse) {

	rate := math.Floor(v.EggGenRate)
	v.EggGenRate = v.EggGenRate + EggFeedHouseProbability(henHouse.Level, feed.Type)
	nowRate := math.Floor(v.EggGenRate)
	number := util.ToInt(nowRate - rate)
	if rate == nowRate {
		//概率没有超过整数不下蛋
	} else {
		//根据概率获取下鸡蛋的个数
		heeAgg(feed, userProperty, v, number, henHouse)
	}

}

//鸡下蛋
func heeAgg(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, number int, henHouse *models.HenHouse) {
	if v.HenType == 1 {
		//免费的鸡下蛋
		for i := 0; i < number; i++ {
			//下彩蛋
			egg := models.Egg{}
			egg.Type = 1
			egg.CreateTimeDay = util.GetCurDayTime()
			egg.UserId = feed.UserId
			egg.HenId = feed.HenId
			egg.CreateTime = util.ToInt(time.Now().Unix())
			//创建下单记录
			models.DB.Create(&egg)
			userProperty.FreeEggs = userProperty.FreeEggs + 1
			//像用户资产表添加彩蛋
			models.DB.Model(&userProperty).Update("FreeEggs")

		}
	} else {
		//乖乖鸡和金鸡
		for i := 0; i < number; i++ {
			//下鸡蛋
			egg := models.Egg{}

			egg.CreateTimeDay = util.GetCurDayTime()
			egg.UserId = feed.UserId
			egg.HenId = feed.HenId
			egg.CreateTime = util.ToInt(time.Now().Unix())
			//根据鸡的种类和鸡舍
			t := rand.Intn(EggGold(v, henHouse))
			if t > 0 && t < 50 {
				egg.Type = 5
				userProperty.GoldEggsFree = userProperty.GoldEggsFree + 1
			} else {
				egg.Type = 2
				userProperty.RealEggs = userProperty.RealEggs + 1
			}
			//创建下蛋记录
			models.DB.Create(&egg)
			//像用户资产表添加鸡蛋
			models.DB.Save(&userProperty)

		}
	}

}

//检查是否还有食料
func CheckFeed(userProperty *models.UserProperty) int {
	if userProperty.NormalFoods > 0 {
		return 1
	} else {
		if userProperty.EnergyFoods > 0 {
			return 2
		} else {
			if userProperty.EnergyFoods > 0 {
				return 3
			} else {
				return 0
			}
		}

	}
}

//根据食料喂食
func DepositFeed(typeFeed int, hen *models.Hen) (flag int) {
	switch typeFeed {
	case 1:
		//系统特制的食料
		flag = 1
		SevenOrOneRecord(1, hen)
		break
	case 2:
		//普通食料
		flag = 2
		SevenOrOneRecord(2, hen)
		break
	case 3:
		//普通食料
		flag = 3
		SevenOrOneRecord(3, hen)
		break
	}
	return
}

//喂养七天，或者不连续喂养
func SevenOrOneRecord(feedType int, hen *models.Hen) {
	//添加喂养动作
	feed := models.Feed{}
	feed.UserId = util.ToInt(hen.UserID)
	//添加昨天的凌晨时间
	feed.CreateTimeDay = util.GetYesDayTime()
	if !models.DB.NewRecord(feed) {
		//判断七天是否连续喂养，然后处理
		SevenRecord(feed, util.ToString(hen.UserID), util.ToString(hen.ID), feedType)

	} else {
		//如果不是则SevenTime置为1
		//创建feed,添加喂养记录
		feedCru := models.Feed{}
		feedCru.CreateTime = util.ToInt(time.Now().Unix())
		feed.CreateTimeDay = util.GetCurDayTime()
		feed.UserId = hen.UserID
		feed.HenId = hen.ID
		feed.Type = feedType
		feed.SevenTime = 1
		models.DB.Create(&feedCru)
		//添加金币
		AddGold(hen.UserID, 1)
		ChangeChilckType(hen.ID, 2)
	}

}

//喂养的七天记录，如果是7天，添加额外的金币，如果不是，那就添加1个金币
func SevenRecord(feed models.Feed, userid, HenId string, feedType int) (flag bool) {
	//如果存在昨天的记录,那就获取累计天数
	models.DB.First(&feed)
	if feed.SevenTime == 7 {
		//创建feed,添加喂养记录
		feedCru := models.Feed{}
		feedCru.CreateTime = util.ToInt(time.Now().Unix())
		feed.CreateTimeDay = util.GetCurDayTime()
		feed.UserId = util.ToInt(userid)
		feed.HenId = util.ToInt(HenId)
		feed.Type = feedType
		feed.SevenTime = 1
		models.DB.Create(&feedCru)
	} else {
		//创建feed,添加喂养记录
		feedCru := models.Feed{}
		feedCru.CreateTime = util.ToInt(time.Now().Unix())
		feed.CreateTimeDay = util.GetCurDayTime()
		feed.UserId = util.ToInt(userid)
		feed.HenId = util.ToInt(HenId)
		feed.Type = feedType
		feed.SevenTime = feed.SevenTime + 1
		models.DB.Create(&feedCru)
		if feed.SevenTime == 7 {
			//连续7天喂养后添加金币
			AddGold(util.ToInt(userid), 1)
		}
	}
	//添加喂养记录后，添加金币
	AddGold(util.ToInt(userid), 1)

	//最后更改鸡的状态
	flag = ChangeChilckType(util.ToInt(HenId), 2)

	return

}
