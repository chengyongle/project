package services

import (
	"myshoppingmall/dao"
	"myshoppingmall/datamodels"
)
import "github.com/jinzhu/gorm"

type ProductService struct {
	productdao dao.Productdao
}
func Newproductservice(db *gorm.DB) *ProductService{
	pd:=dao.Productdao{db}
	return  &ProductService{pd}
}
func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productdao.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return  p.productdao.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return  p.productdao.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return  p.productdao.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productdao.Update(product)
}
func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productdao.SubProductNum(productID)
}
