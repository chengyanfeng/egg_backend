package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type Gorm struct {
}

func (User) TableName() string {

	return "user"
}
func (HenHouse) TableName() string {

	return "henhouse"
}
func (Hen) TableName() string {

	return "hen"
}
func (Egg) TableName() string {

	return "egg"
}

//数据库初始化
func init() {
	var err error
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		Mysqlconn["username"],
		Mysqlconn["password"],
		Mysqlconn["host"],
		Mysqlconn["port"],
		Mysqlconn["name"],
	)
	fmt.Print("mysqlconn:")
	fmt.Println(Mysqlconn)
	DB, err = gorm.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	if DB.HasTable("user") {
		//自动添加模式
		DB.AutoMigrate(&User{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&User{})
	}
	if DB.HasTable("henhouse") {
		//自动添加模式
		DB.AutoMigrate(&HenHouse{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&HenHouse{})
	}
	if DB.HasTable("hen") {
		//自动添加模式
		DB.AutoMigrate(&Hen{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&Hen{})
	}
	if DB.HasTable("egg") {
		//自动添加模式
		DB.AutoMigrate(&Egg{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&Egg{})
	}
}
