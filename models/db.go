package models

type P map[string]interface{}

var Mysqlconn = P{

	"username": "root",
	"password": "12345678",
	"host":     "localhost",
	"port":     3306,
	"name":     "test",
}
