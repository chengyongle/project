package services

import (
	"github.com/jinzhu/gorm"
	"myshoppingmall/dao"
	"myshoppingmall/datamodels"

)

type OrderService struct {
	orderdao dao.Orderdao
}
func NeworderService(db *gorm.DB) *OrderService{
	od:=dao.Orderdao{db}
	return  &OrderService{od}
}
//调用dao层函数
func (o *OrderService) GetOrderByID(orderID int64) (*datamodels.Order, error) {
	return o.orderdao.SelectByKey(orderID)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return  o.orderdao.SelectAll()
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return  o.orderdao.Delete(orderID)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return  o.orderdao.Insert(order)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderdao.Update(order)
}
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (int64 ,error) {
	order :=&datamodels.Order{
		UserID:      message.UserID,
		ProductID: message.ProductID,
		OrderStatus: 0,
	}
	return o.orderdao.Insert(order)
}