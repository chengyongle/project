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

type ProductController struct {
	ps *services.ProductService
}

func NewProductController(db *gorm.DB)*ProductController{
	return &ProductController{services.Newproductservice(db)}
}

func (p *ProductController) Router(app *gin.Engine) {
	app.GET("/product/all", p.Getall)
	app.POST("/product/update",p.Update)
	app.GET("/product/add",p.Gadd)
	app.POST("/product/add",p.Padd)
	app.GET("/product/manager",p.Manager)
	app.GET("/product/delete",p.Delete)
}
//获取全部
func (p *ProductController) Getall(c *gin.Context){
	productArray, _ := p.ps.GetAllProduct()
	c.HTML(http.StatusOK,"product/view.html",gin.H{"productArray": productArray})
}

//更新
func (p *ProductController) Update(c *gin.Context) {
	product := datamodels.Product{}
	c.ShouldBind(&product)
	err := p.ps.UpdateProduct(&product)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}
//增加
func (p *ProductController) Gadd(c *gin.Context){
	c.HTML(http.StatusOK,"product/add.html",nil)
}

func (p *ProductController) Padd(c *gin.Context) {
	product := datamodels.Product{}
	//参数绑定
	c.ShouldBind(&product)
	_, err := p.ps.InsertProduct(&product)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}

func (p *ProductController) Manager(c *gin.Context)  {
	idstring:=c.Query("id")
	id, _ := strconv.ParseInt(idstring, 10, 64)
	product, err := p.ps.GetProductByID(id)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.HTML(http.StatusOK,"product/manager.html",gin.H{"product": product})
}

func (p *ProductController) Delete(c *gin.Context) {
	idstring:=c.Query("id")
	id, _ := strconv.ParseInt(idstring, 10, 64)

	isOk := p.ps.DeleteProductByID(id)
	if !isOk {
		fmt.Printf("删除商品失败，ID为：" + idstring)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"all")
}
