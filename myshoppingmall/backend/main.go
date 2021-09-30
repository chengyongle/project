package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"myshoppingmall/backend/web/controllers"
	"myshoppingmall/datamodels"
)


func main() {
	db:=sqlinit()
	defer db.Close()
	app:=gin.Default()
	app.LoadHTMLGlob("backend/web/views/**/*")
	registerRouter(app,db)
	app.Run()

}
//注册路由
func registerRouter(router *gin.Engine,db *gorm.DB) {
	controllers.NewProductController(db).Router(router)
	controllers.NewOrderController(db).Router(router)
}
//获取数据库连接
func sqlinit()*gorm.DB{
	dst:="root:123456@tcp(127.0.0.1:3306)/database1?charset=utf8"
	db, err := gorm.Open("mysql", dst)
	if err!=nil{
		fmt.Printf("err:%v",err)
		return nil
	}

	if table := db.HasTable(datamodels.Product{});!table {
		db.CreateTable(datamodels.Product{})
	}
	if table := db.HasTable(datamodels.Order{});!table {
		db.CreateTable(datamodels.Order{})
	}
	return  db
}
