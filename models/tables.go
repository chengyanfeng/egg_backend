package models

import (
	"time"
)

type Model struct {
}

//用户表
type User struct {
	ID           int `gorm:"primary_key"`
	CreateTime   int `gorm:"CreateTime"`
	LastLogIn    int `gorm:"LastLogIn"`
	DeletedAt    *time.Time
	LoginTimes   int    `gorm:"LoginTimes"`
	AddressInfo  string `gorm:"column:AddressInfo"`
	WXOpenID     string `gorm:"column:WXOpenID"`
	WXNickName   string `gorm:"column:WXNickName"`
	WXHeadImg    string `gorm:"column:WXHeadImg"`
	Mobile       string `gorm:"column:Mobile"`
	PwdHash      string `gorm:"column:PwdHash"`
	LastLogInDay int    `gorm:"column:LastLogInDay"`
	SevenTimeSum int    `gorm:"column:SevenTimeSum"`
	PropertyID   int    `gorm:"column:PropertyID"`
}

//充值表
type ChargeOrder struct {
	ID         int `gorm:"primary_key"`
	UserID     int `gorm:"UserID"` //userID 外键
	Amount     int `gorm:"Amount"` //充值额（RMB分）
	CreateTime time.Time
	PayID      string `gorm:"PayID"`    //支付单号（取决于微信或支付宝）
	PayState   int    `gorm:"PayState"` //支付状态（0：未支付，1：已支付，2:支付失败）
}

//用户资产表
type UserProperty struct {
	ID                int     `gorm:"primary_key"`              //id
	UserID            int     `gorm:"UserID"`                   //用户ID （关联用户表，1对1）
	Coins             int     `gorm:"column:Coins"`             //金币数量
	Wallet            string  `gorm:"column:Wallet"`            //BSTK Wallet
	RealEggs          int     `gorm:"column:RealEggs"`          //鸡蛋数量
	FreeEggs          int     `gorm:"column:FreeEggs"`          //彩蛋数量
	GoldEggsFree      int     `gorm:"column:GoldEggsFree"`      //可操作金蛋数量
	GoldEggMarketLock int     `gorm:"column:GoldEggMarketLock"` //挂在市场上的数量
	HenHouseID        float64 `gorm:"column:HenHouseID"`        //鸡舍ID（关联鸡舍表，1对1）
	Hens              int     `gorm:"column:Hens"`              //所拥有鸡ID列表,此项属性低频变化，用户购买，或者孵化，或者鸡死亡触发变化，平常避免联合查询
	NormalFoods       int     `gorm:"column:NormalFoods"`       //普通饲料数量
	FastFoods         int     `gorm:"column:FastFoods"`         //快速成长饲料数量
	HasDag            int     `gorm:"column:HasDag"`            //是否有狗 (0, 1)
}

//商品表
type Shop struct {
	ID        int    `gorm:"primary_key"`      //商品ID
	Name      string `gorm:"column:Name"`      //商品名称
	Type      int    `gorm:"column:Type"`      //商品类型（1:乖乖鸡，2:普通饲料，3:能量饲料，4:快速饲料..., 类型通过和客户端共享JSON配置文件）
	PriceCent int    `gorm:"column:PriceCent"` //金币价格（以分为单位，1/100）
	Amount    int    `gorm:"column:Amount"`    //库存
}

//商品订单表
type ShopOrder struct {
	ID         int `gorm:"primary_key"`
	UserID     int `gorm:"column:UserID"`
	CreateTime int `gorm:"column:CreateTime"` //订单创建时间 UTC
	ShopID     int `gorm:"column:ShopID"`     //商品ID
	Amount     int `gorm:"column:Amount"`     //购买数量
	Coins      int `gorm:"column:Coins"`      //支付总额（金币分）
	State      int `gorm:"column:State"`      //0: 待支付，1: 支付成功
}

//系统回收鸡蛋订单表
type EggWithdrawOrder struct {
	ID         int `gorm:"primary_key"`
	UserID     int `gorm:"column:UserID"`
	CreateTime int `gorm:"column:CreateTime"` //创建时间戳
	Amount     int `gorm:"column:Amount"`     //数量
	PriceCent  int `gorm:"column:PriceCent"`  //单价（金币分）
}

//用户提现鸡蛋订单表
type EggTakenOrder struct {
	ID          int    `gorm:"primary_key"`
	UserID      int    `gorm:"column:UserID"`      //用户ID
	CreateTime  int    `gorm:"column:CreateTime"`  //创建时间戳
	Amount      int    `gorm:"column:Amount"`      //提现鸡蛋数量
	Address     string `gorm:"column:Address"`     //用户收货信息，从用户表中获得
	DeliverInfo string `gorm:"column:DeliverInfo"` //快递信息
	State       int    `gorm:"column:State"`       //订单状态
}

//鸡表
type Hen struct {
	ID          int     `gorm:"primary_key"`        //id
	RealTag     string  `gorm:"column:RealTag"`     //实体鸡脚环唯一标识
	Name        string  `gorm:"column:Name"`        //鸡昵称
	CreateTime  int     `gorm:"column:CreateTime"`  //创建时间，鸡生日（UTC， timestamp)
	State       int     `gorm:"column:State"`       //当前状态：（1:饥饿，2:吃饱，3:无人看管，4:出游）
	HenType     int     `gorm:"column:HenType"`     //鸡类型：（1:免费鸡，2:乖乖鸡，3:金鸡，4:鸡雏，5:公鸡）
	EggType     int     `gorm:"column:EggType"`     //产蛋类型：（0:无法产蛋，1:彩蛋，2:鸡蛋），都有概率产出金蛋（3）
	LifeTime    int     `gorm:"column:LifeTime"`    //鸡龄: 蛋鸡365天，鸡雏喂养30天转为乖乖鸡，可通过道具加速
	Lifes       int     `gorm:"column:Lifes"`       //生命数，3条
	LifeValue   int     `gorm:"column:LifeValue"`   //生命值，一天不喂养（喂养以天为单位，0-24点任意时刻），进入生命值倒计时（72小时），State转为1，倒计时内喂养，解除倒计时，倒计时到减除1条Life直至完全死亡
	EggGenRate  float64 `gorm:"column:EggGenRate"`  //产蛋率（每日产蛋数量，default: 0.667，即3天2枚，可转化为36小时1枚)，针对彩蛋，鸡蛋，金蛋
	GoldEggRate int     `gorm:"column:GoldEggRate"` //产蛋为金蛋概率，每产出一枚蛋时，依据此概率进行金蛋转换(千分之）
	Skins       string  `gorm:"column:Skins"`       //当前使用的道具列表
	HenHouseID  int     `gorm:"column:HenHouseID"`  //所属鸡舍 （关联鸡舍表，多对1）
	UserID      int     `gorm:"column:UserID"`      //所属用户 （关联用户表，多对1）
	Deposit     int     `gorm:"column:Deposit"`     //0，未托管，1托管。
}

//鸡舍表
type HenHouse struct {
	ID         int    `gorm:"primary_key"`       //id
	Level      int    `gorm:"column:Level"`      //等级
	Tools      string `gorm:"column:Tools"`      //道具列表，用户购买的鸡舍道具列表
	CleanState int    `gorm:"column:CleanState"` //	清洁程度
	UserID     int    `gorm:"column:UserID"`     //所属用户 (关联用户表，1对1）
}

//喂养表
type Feed struct {
	ID            int `gorm:"primary_key"`          //ID
	Type          int `gorm:"column:Type"`          //喂养的食料的种类（1:乖乖鸡，2:普通饲料，3:能量饲料，4:快速饲料..., 类型通过和客户端共享JSON配置文件）
	CreateTime    int `gorm:"column:CreateTime"`    //喂养的时间
	UserId        int `gorm:"column:UserId"`        //用户ID
	HenId         int `gorm:"column:HenId"`         //喂养的鸡的id
	CreateTimeDay int `gorm:"column:CreateTimeDay"` //喂养的当天的零点时间，为了方便。
	SevenTime     int `gorm:"column:SevenTime"`     //七天连续的天数
}
type Egg struct {
	ID            int `gorm:"primary_key"`          //ID
	Type          int `gorm:"column:Type"`          //鸡蛋的类型
	CreateTime    int `gorm:"column:CreateTime"`    //创建时间
	UserId        int `gorm:"column:UserId"`        //用户ID
	HenId         int `gorm:"column:HenId"`         //鸡的ID
	CreateTimeDay int `gorm:"column:CreateTimeDay"` //下蛋当天的零点时间，为了方便。
	Source        int `gorm:"column:Source"`        //鸡蛋的来源 1.小鸡下的，2.兑换的。
	Sell          int `gorm:"column:Sell"`          //是否售卖，或者兑换金币，或其他蛋种
}

//市场表
type Market struct {
	Id            int `gorm:"ID"`
	Seller        int `gorm:"Seller"`     //卖家UserID
	Buyer         int `gorm:"Buyer"`      //买家UserID 未成交为null
	CreateTime    int `gorm:"CreateTime"` //创建时间戳
	UpdateTime    int `gorm:"UpdateTime"` //更新时间戳（例如卖家更新价格，更新金蛋数量）
	DealTime      int `gorm:"DealTime"`   //成交时间戳
	State         int `gorm:"State"`      //0:待售，1:成交
	Type          int `gorm:"Type"`       //交易物品类型（1:乖乖鸡，2:金鸡，3:金蛋）
	ItemID        int //对应的实际物品ID，目前主要是Hen ID, 如果Type是金蛋，此处对应的是UserPropertyID, UserProperty--GoldEggMarketLock 确定了售卖数量
	ItemPriceCent int `gorm:"ItemPriceCent"` //单价（金币分）
}

//配置表
type Config struct {
	EggGenHour                   int `gorm:"EggGenHour"`                   //产蛋所需小时数，default:36
	HenLifeHours                 int `gorm:"HenLifeHours"`                 //鸡龄 default: 365 x 24
	ChickenNormalLifeHours       int `gorm:"ChickenNormalLifeHours"`       //雏鸡正常喂养时间 30 x 24
	ChickenFastLifeHours         int `gorm:"ChickenFastLifeHours"`         //使用加速包后的喂养时间 3 x 24
	FreeEggConvertRate           int `gorm:"FreeEggConvertRate"`           //彩蛋兑换鸡蛋比率 3
	RealEggSellMin               int `gorm:"RealEggSellMin"`               //系统回收鸡蛋最小数量 10
	RealEggSellPrice             int `gorm:"RealEggSellPrice"`             //系统回收鸡蛋价格 （金币分） 待定
	RealEggTakenMin              int `gorm:"RealEggTakenMin"`              //体现鸡蛋最小数量 30
	GoldEggTransferChickenRate   int `gorm:"GoldEggTransferChickenRate"`   //金蛋孵化小鸡概率（千分之） 300-500
	GoldEggTransferSmallBstkRate int `gorm:"GoldEggTransferSmallBstkRate"` //金蛋孵化小额BSTK概率（千分之）500-700
	GoldEggTransferHugeBstkRate  int `gorm:"GoldEggTransferHugeBstkRate"`  //金蛋孵化大额BSTK概率（千分之）待定 （3项概率和为1000）
	EnergeFoodAffection          int `gorm:"EnergeFoodAffection"`          //能量饲料产蛋影响倍数（只影响第二天）x2
	HenVisitHours                int `gorm:"HenVisitHours"`                //鸡外出串门小时数 1-24小时 随机
}

//鸡舍配置表
type HenHouseConfig struct {
	Level            int `gorm:"Level"`            //生产等级
	GoldEggRateAdd   int `gorm:"GoldEggRateAdd"`   //生产时间
	EggGenRateAdd    int `gorm:"EggGenRateAdd"`    //(千分之)
	UpgradePriceCent int `gorm:"UpgradePriceCent"` //(金币分)
}
