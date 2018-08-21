package router

import (
	"github.com/gin-gonic/gin"
	"Egg/controller"
)
var Defaultrouter = gin.Default()    //获得路由实例
func init() {


	//注册接口demo
	//添加中间件
	Defaultrouter.Use(controller.Middleware)
	Defaultrouter.GET("/simple/server/get", controller.GetHandler)
	Defaultrouter.POST("/simple/server/post", controller.PostHandler)
	Defaultrouter.PUT("/simple/server/put", controller.PutHandler)
	Defaultrouter.DELETE("/simple/server/delete", controller.DeleteHandler)


	/***********************************----------------路由拦截器----------------******************************/
	//添加中间件
	Defaultrouter.Use(controller.LoginFilter)
	/***********************************---------------以下为正式路由-----------------******************************/
	//支付宝支付接口
	Defaultrouter.GET("/zhifubao/pay", controller.Zhifubaopay)
	//支付宝回掉函数
	Defaultrouter.POST("/zhifubao/return", controller.Return)

	//微信登陆接口
	Defaultrouter.GET("/wxlogin/url", controller.RedirecturlHandler)
	//微信返回跳转接口
	Defaultrouter.GET("/wxlogin/index", controller.IndexHandler)

}
