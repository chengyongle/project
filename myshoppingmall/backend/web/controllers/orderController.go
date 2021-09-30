package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"myshoppingmall/datamodels"
	"myshoppingmall/services"
	"net/http"
	"strconv"
)

type OrderController struct {
	os *services.OrderService
}

func NewOrderController(db *gorm.DB)*OrderController{
	return &OrderController{services.NeworderService(db)}
}

func (o *OrderController) Router(app *gin.Engine) {
	app.GET("/order/all", o.Getall)
	app.POST("/order/update",o.Update)
	app.GET("/order/add",o.Gadd)
	app.POST("/order/add",o.Padd)
	app.GET("/order/manager",o.Manager)
	app.GET("/order/delete",o.Delete)
}

func (o *OrderController) Getall(c *gin.Context){
	orderArray, _ := o.os.GetAllOrder()
	c.HTML(http.StatusOK,"order/view.html",gin.H{"orderArray": orderArray})
}


func (o *OrderController) Update(c *gin.Context) {
	order := datamodels.Order{}
	c.ShouldBind(&order)
	err := o.os.UpdateOrder(&order)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}
func (o *OrderController) Gadd(c *gin.Context){
	c.HTML(http.StatusOK,"order/add.html",nil)
}

func (o *OrderController) Padd(c *gin.Context) {
	order := datamodels.Order{}
	c.ShouldBind(&order)
	_, err := o.os.InsertOrder(&order)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}

func (o *OrderController) Manager(c *gin.Context)  {
	idstring:=c.Query("id")
	id, _ := strconv.ParseInt(idstring, 10, 64)
	order, err := o.os.GetOrderByID(id)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.HTML(http.StatusOK,"order/manager.html",gin.H{"order": order})
}

func (o *OrderController) Delete(c *gin.Context) {
	idstring:=c.Query("id")
	id, _ := strconv.ParseInt(idstring, 10, 64)

	isOk := o.os.DeleteOrderByID(id)
	if !isOk {
		fmt.Printf("删除商品失败，ID为：" + idstring)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}
