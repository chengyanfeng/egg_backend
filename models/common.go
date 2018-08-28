package models

import (
	"egg_backend/util"
	"math"
	"math/rand"
	"time"
)

//添加金币
func AddGold(userId, number int) {
	propertyEvenDay := UserProperty{}
	propertyEvenDay.UserID = userId
	DB.First(&propertyEvenDay)
	propertyEvenDay.Coins = propertyEvenDay.Coins + number
	DB.Model(&propertyEvenDay).Update("coins")
}

//更改鸡的状态
func ChangeChilckType(HenId, Type int) bool {
	if !DB.NewRecord(HenId) {
		Hen := Hen{}
		Hen.ID = HenId
		DB.First(&Hen)
		Hen.State = Type
		DB.Model(&Hen).Update("state")
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
func EggGold(v *Hen, henHouse *HenHouse) (randint int) {
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
	if v.HenType == 3 {
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
func FreeHee(feed *Feed, userProperty *UserProperty, v *Hen, henHouse *HenHouse) {
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
func DarlingHee(feed *Feed, userProperty *UserProperty, v *Hen, henHouse *HenHouse) {

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
func GoldHee(feed *Feed, userProperty *UserProperty, v *Hen, henHouse *HenHouse) {

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

//获取雏鸡，如果是雏鸡那么不生蛋
func LitileHee(feed *Feed, userProperty *UserProperty, v *Hen, henHouse *HenHouse) {
	//根据概率获取下鸡蛋的个数
	heeAgg(feed, userProperty, v, 0, henHouse)
}

//鸡下蛋
func heeAgg(feed *Feed, userProperty *UserProperty, v *Hen, number int, henHouse *HenHouse) {
	if v.HenType == 1 {
		//免费的鸡下蛋
		for i := 0; i < number; i++ {
			//下彩蛋
			egg := Egg{}
			egg.Type = 1
			egg.CreateTimeDay = util.GetCurDayTime()
			egg.UserId = feed.UserId
			egg.HenId = feed.HenId
			egg.CreateTime = util.ToInt(time.Now().Unix())
			//创建下单记录
			DB.Create(&egg)
			userProperty.FreeEggs = userProperty.FreeEggs + 1
			//像用户资产表添加彩蛋
			DB.Model(&userProperty).Update("FreeEggs")

		}
	} else if v.HenType == 4 {
		feedcount := Feed{}
		count := 0
		//如果是雏鸡，那么不下蛋，查看累计喂的食料
		DB.Model(&feedcount).Where("UserId=? AND HenId=?", v.UserID, v.ID).Count(&count)
		if count > 30 {
			//雏鸡转换为乖乖鸡
			v.HenType = 2
			DB.Model(v).Update("hen_type")
		}
	} else {
		//乖乖鸡和金鸡
		for i := 0; i < number; i++ {
			//下鸡蛋
			egg := Egg{}

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
			DB.Create(&egg)
			//像用户资产表添加鸡蛋
			DB.Save(&userProperty)

		}
	}

}

//检查是否还有食料
func CheckFeed(userProperty *UserProperty) int {
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
func DepositFeed(typeFeed int, hen *Hen) (flag int) {
	switch typeFeed {
	case 1:
		//系统赠送的食料
		flag = 1
		SevenOrOneRecord(1, hen)
		break
	case 2:
		//普通食料
		flag = 2
		SevenOrOneRecord(2, hen)
		break
	case 3:
		//快速食料
		flag = 3
		SevenOrOneRecord(3, hen)
		break
	}
	return
}

//喂养七天，或者不连续喂养
func SevenOrOneRecord(feedType int, hen *Hen) {
	//通过食料类型喂养
	FeedType(feedType, hen.UserID, hen.ID)
}

//判断喂养的饲料
func FeedType(feedType int, userId, henId int) {
	//通过类型像道具中获取持续天数
	foodType := FoodStuff{}
	foodType.Type = feedType
	DB.First(&foodType)
	for i := 0; i < foodType.ContinueDay; i++ {
		//创建喂养记录
		feed := Feed{}
		feed.CreateTimeDay = util.GetCurDayTime() + i*86400
		feed.CreateTime = util.ToInt(time.Now().Unix())
		feed.ContinueDay = (foodType.ContinueDay - i)
		feed.Type = feedType
		feed.UserId = userId
		feed.HenId = henId
		//获取前一天的喂养记录
		YesTaday := Feed{}
		YesTaday.HenId = henId
		YesTaday.UserId = userId
		YesTaday.CreateTimeDay = util.GetCurDayTime() + (i-1)*86400
		if !DB.NewRecord(YesTaday) {
			feed.SevenTime = (i % 7) + 1
		} else {
			feed.SevenTime = 1
		}
		//查询是否有今天的记录
		goldFeed := Feed{}
		count := 0
		DB.Model(&goldFeed).Where("UserId=? AND HenId=? AND CreateTimeDay=?", userId, henId, util.ToInt(util.GetCurDayTime())).Count(&count)
		if count > 0 {
			//说明当天已经喂过了，就不再添加金币了。
			//创建记录
			DB.Create(&feed)

		} else {
			//创建记录
			DB.Create(&feed)
			//添加金币
			AddGold(userId, 1)
			if feed.SevenTime == 7 {
				//添加额外的金币
				AddGold(userId, 1)
			}
		}

	}

}
