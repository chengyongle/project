package datamodels

type Order struct {
	ID          int64 `json:"id"  gorm:"column:id;AUTO_INCREMENT"`
	UserID      int64 `json:"Username"  gorm:"column:userid"`
	ProductID   int64 `json:"Productname"  gorm:"column:productid"`
	OrderStatus int `json:"OrderStatus"  gorm:"column:orderstatus"`
}

const (
	OrderWait    = iota
	OrderSuccess //1
	OrderFailed  //2
)
