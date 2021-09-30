package main
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"myshoppingmall/datamodels"
	"myshoppingmall/fronted/web/controllers"
	"myshoppingmall/rabbitmq"
)

func main() {
	db:=sqlinit()
	defer  db.Close()
	app:=gin.Default()
	app.LoadHTMLGlob("fronted/web/views/**/*")
	rb:=rabbitmq.NewRabbitMQSimple("myProduct")
	registerRouter(app,db,rb)
	app.Run()

}
func registerRouter(router *gin.Engine,db *gorm.DB,rb *rabbitmq.RabbitMQ) {
	controllers.NewUserController(db).Router(router)
	controllers.NewProductController(db,rb).Router(router)
}
//获取数据库连接
func sqlinit()*gorm.DB{
	dst:="root:123456@tcp(127.0.0.1:3306)/database1?charset=utf8"
	db, err := gorm.Open("mysql", dst)
	if err!=nil{
		fmt.Printf("err:%v",err)
		return nil
	}
	if table := db.HasTable(datamodels.User{});!table {
		db.CreateTable(datamodels.User{})
	}
	return  db
}
