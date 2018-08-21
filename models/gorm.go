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
		DB.AutoMigrate(&User{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&User{})
	}

}
