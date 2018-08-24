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

//鸡的生蛋概率
func EggProbability() float64 {
	Probability := []float64{0.6, 0.7, 0.8, 0.9, 1}
	t := rand.Intn(4)
	return Probability[t]
}

//获取免费鸡生蛋跟据喂养的饲料，根据概率，产生的蛋
func FreeHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen) {
	//如果食料为系统赠送的饲料
	if feed.Type == 1 {
		freehee(feed, userProperty, v, 1)
	}

	//如果是普通的饲料
	if feed.Type == 2 {
		rate := math.Floor(v.EggGenRate)
		v.EggGenRate = v.EggGenRate + EggProbability()
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}
	//如果加速包
	if feed.Type == 3 {
		rate := math.Floor(v.EggGenRate)
		//概率目前暂定X2
		v.EggGenRate = v.EggGenRate + EggProbability()*2
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}

}

//免费公鸡下蛋
func freehee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, number int) {
	for i := 0; i < number; i++ {
		//下彩蛋
		egg := models.Egg{}
		egg.Type = 1
		egg.CreateTimeDay = util.GetCurDayTime()
		egg.UserId = feed.UserId
		egg.HenId = feed.HenId
		egg.CreateTime = util.ToInt(time.Now().Unix())
		models.DB.Save(&egg)
		userProperty.FreeEggs = userProperty.FreeEggs + 1
		//像用户资产表添加彩蛋
		models.DB.Model(&userProperty).Update("FreeEggs")

	}

}

//获取乖乖鸡生蛋跟据喂养的饲料，根据概率，产生的蛋
func DarlingHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen) {
	//如果食料为系统赠送的饲料
	if feed.Type == 1 {
		freehee(feed, userProperty, v, 1)
	}

	//如果是普通的饲料
	if feed.Type == 2 {
		rate := math.Floor(v.EggGenRate)
		v.EggGenRate = v.EggGenRate + EggProbability()
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}
	//如果加速包
	if feed.Type == 3 {
		rate := math.Floor(v.EggGenRate)
		//概率目前暂定X2
		v.EggGenRate = v.EggGenRate + EggProbability()*2
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}

}

//免费乖乖鸡下蛋
func darlingHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, number int) {
	for i := 0; i < number; i++ {
		//下彩蛋
		egg := models.Egg{}
		egg.Type = 1
		egg.CreateTimeDay = util.GetCurDayTime()
		egg.UserId = feed.UserId
		egg.HenId = feed.HenId
		egg.CreateTime = util.ToInt(time.Now().Unix())
		models.DB.Save(&egg)
		userProperty.FreeEggs = userProperty.FreeEggs + 1
		//像用户资产表添加彩蛋
		models.DB.Model(&userProperty).Update("FreeEggs")

	}

}

//获取金鸡生蛋跟据喂养的饲料，根据概率，产生的蛋
func GoldHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen) {
	//如果食料为系统赠送的饲料
	if feed.Type == 1 {
		freehee(feed, userProperty, v, 1)
	}

	//如果是普通的饲料
	if feed.Type == 2 {
		rate := math.Floor(v.EggGenRate)
		v.EggGenRate = v.EggGenRate + EggProbability()
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}
	//如果加速包
	if feed.Type == 3 {
		rate := math.Floor(v.EggGenRate)
		//概率目前暂定X2
		v.EggGenRate = v.EggGenRate + EggProbability()*2
		nowRate := math.Floor(v.EggGenRate)
		number := util.ToInt(nowRate - rate)
		if rate == nowRate {
			//概率没有超过整数不下蛋
		} else {
			//根据概率获取下鸡蛋的个数
			freehee(feed, userProperty, v, number)
		}
	}

}

//金鸡下蛋
func goldHee(feed *models.Feed, userProperty *models.UserProperty, v *models.Hen, number int) {
	for i := 0; i < number; i++ {
		//下彩蛋
		egg := models.Egg{}
		egg.Type = 1
		egg.CreateTimeDay = util.GetCurDayTime()
		egg.UserId = feed.UserId
		egg.HenId = feed.HenId
		egg.CreateTime = util.ToInt(time.Now().Unix())
		models.DB.Save(&egg)
		userProperty.FreeEggs = userProperty.FreeEggs + 1
		//像用户资产表添加彩蛋
		models.DB.Model(&userProperty).Update("FreeEggs")

	}

}
