package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"myshoppingmall/datamodels"
	"myshoppingmall/rabbitmq"
	"myshoppingmall/services"
)

//rabbitmq消费者
func main() {
	db:=sqlinit()
	defer db.Close()
	os:=*services.NeworderService(db)
	ps:=*services.Newproductservice(db)
	rabbitmqConsumeSimple :=rabbitmq.NewRabbitMQSimple("myProduct")
	rabbitmqConsumeSimple.ConsumeSimple(os,ps)//传入订单和库存server
}
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
