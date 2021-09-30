package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"myshoppingmall/datamodels"
	"myshoppingmall/rabbitmq"
	"myshoppingmall/services"
	"net/http"
	"strconv"
)

type ProductController struct {
	ps       *services.ProductService
	os       *services.OrderService
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewProductController(db *gorm.DB,rb *rabbitmq.RabbitMQ) *ProductController {
	return &ProductController{services.Newproductservice(db), services.NeworderService(db),rb}
}

func (p *ProductController) Router(app *gin.Engine) {
	app.GET("/product/detail", p.GetDetail)
	app.GET("/product/order", p.GetOrder)
	app.GET("/product/all", p.Getall)
}
func (p *ProductController) Getall(c *gin.Context) {
	productArray, _ := p.ps.GetAllProduct()
	c.HTML(http.StatusOK, "product/view.html", gin.H{"productArray": productArray})
}
func (p *ProductController) GetDetail(c *gin.Context) {
	pidstring := c.Query("id")
	pid, _ := strconv.ParseInt(pidstring, 10, 64)
	product, err := p.ps.GetProductByID(pid)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	c.HTML(http.StatusOK, "product/view.html", gin.H{"product": product})
}

func (p *ProductController) GetOrder(c *gin.Context){
	productString := c.Query("productID")
	userString, err := c.Cookie("uid")
	if err != nil {
		//未登录
		c.HTML(http.StatusOK, "user/login.html", nil)
	}
	productID, _ := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	//创建消息体
	message := datamodels.NewMessage(userID, productID)
	//类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		log.Fatal(err)
	}

	c.String(200,"true")
}
