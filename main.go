package main

import (
	"egg_backend/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

var routerDefaul = router.Defaultrouter

func main() {

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode

	//监听端口
	http.ListenAndServe(":8005", routerDefaul)
}
