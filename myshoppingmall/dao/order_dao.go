package dao

import (
	"github.com/jinzhu/gorm"
	"myshoppingmall/datamodels"
)

type Orderdao struct {
	Db  *gorm.DB
}
func (o *Orderdao) Insert(order *datamodels.Order) (int64,error) {
	err:=o.Db.Create(order).Error
	if err!=nil{
		return  0,err
	}
	var id []int64
	o.Db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return  id[0],err
}

func (o *Orderdao) Delete(orderID int64) bool{
	pro:=&datamodels.Product{
		ID: orderID,
	}
	err:=o.Db.Delete(pro).Error
	if err!=nil{
		return  false
	}
	return  true
}

func (o *Orderdao) Update(order *datamodels.Order) error {
	err:=o.Db.Model(order).Updates(order).Error
	return err
}

func (o *Orderdao) SelectByKey(orderID int64) (*datamodels.Order, error) {
	orderResult := &datamodels.Order{}
	if err:=o.Db.Where("id=?",orderID).First(orderResult).Error;
		err==gorm.ErrRecordNotFound{
		return &datamodels.Order{},gorm.ErrRecordNotFound
	}
	return  orderResult,nil
}

func (o *Orderdao) SelectAll() (orderArray []*datamodels.Order, err error) {
	orders:=[]datamodels.Order{}
	if err:=o.Db.Find(&orders).Error;err!=nil{
		return []*datamodels.Order{},err
	}
	res:=make([]*datamodels.Order,len(orders))
	for i:=range res{
		res[i]=&orders[i]
	}
	return  res,nil
}