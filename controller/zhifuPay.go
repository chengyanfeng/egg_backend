package controller

import (
	"fmt"
	"github.com/smartwalle/alipay"
	"github.com/gin-gonic/gin"
	"egg_backend/def"
	"net/http"
)


var client = alipay.New(def.ZHIFUBAOAPPID, "2088802940812132", def.ZHIFUBAO_KEY, def.ZHIFUBAOprivateKey, false)

//支付信息
func  Zhifubaopay(c *gin.Context) {
	var p = alipay.AliPayTradeWapPay{}
	//回调函数
	p.NotifyURL = "http://192.144.176.213:8070/Return"
	p.ReturnURL = "http://192.144.176.213:8070/apliy"
	p.Subject = "这是测试"
	p.OutTradeNo = "2342341233121w3q2eq131w2"
	p.TotalAmount = "10.00"
	p.ProductCode = "商品编码"

	var html, _ = client.TradeWapPay(p)
	c.JSON(http.StatusOK, []byte(fmt.Sprintf("get success! %s\n", html)))
}
//回调函数，支付成功与否，支付宝主动通知这个方法，如果不返回sucess那支付宝会按时发查询
func Return(c *gin.Context){
	req1 :=c.Request
	fmt.Print(req1,"--------------req1-------------")
	fmt.Print(req1.Form,"------------req1.Form----------")

	ok, err := client.VerifySign(req1.Form)
	//如果校验成功后返回success
	fmt.Println(ok, err)
	c.JSON(http.StatusOK, []byte("success"))

}
