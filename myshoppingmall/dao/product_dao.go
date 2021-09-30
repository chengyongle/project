package dao

import (
	"github.com/jinzhu/gorm"
	"myshoppingmall/datamodels"
	"strconv"
)

type Productdao struct {
	Db  *gorm.DB
}


func (p *Productdao) Insert(product *datamodels.Product) (int64, error) {
	err:=p.Db.Create(product).Error
	if err!=nil{
		return  0,err
	}
	var id []int64
	p.Db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return  id[0],err
}

func (p *Productdao) Delete(productID int64) bool {
	pro:=&datamodels.Product{
		ID: productID,
	}
	err:=p.Db.Delete(pro).Error
	if err!=nil{
		return  false
	}
	return  true
}

func (p *Productdao) Update(product *datamodels.Product) error {
	err:=p.Db.Model(product).Updates(product).Error
	return err
}

func (p *Productdao) SelectByKey(productID int64) (*datamodels.Product, error) {
	productResult := &datamodels.Product{}
	if err:=p.Db.Where("id=?",productID).First(productResult).Error;
		err==gorm.ErrRecordNotFound{
		return &datamodels.Product{},gorm.ErrRecordNotFound
	}
	return  productResult,nil
}

func (p *Productdao) SelectAll() ([]*datamodels.Product, error) {
	products:=[]datamodels.Product{}
	if err:=p.Db.Find(&products).Error;err!=nil{
		return []*datamodels.Product{},err
	}
	res:=make([]*datamodels.Product,len(products))
	for i:=range res{
		res[i]=&products[i]
	}
	return  res,nil
}
func (p *Productdao) SubProductNum(productID int64) error {
	id:=strconv.FormatInt(productID,10)
	err:=p.Db.Model(&datamodels.Product{}).Where("id=?",id).Update("productnum", gorm.Expr("productnum-?", 1)).Error
	return err
}


