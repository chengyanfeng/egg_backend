package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

var DB *gorm.DB

type Gorm struct {
}

func (UserInfo) TableName() string {

	return "users"
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
	if DB.HasTable("users") {
		//自动添加模式
		DB.AutoMigrate(&UserInfo{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&UserInfo{})
	}

}
