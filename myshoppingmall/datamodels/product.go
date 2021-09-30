package datamodels

type Product struct {
	ID           int64  `json:"id"  gorm:"column:id;AUTO_INCREMENT"`
	ProductName  string `json:"ProductName" gorm:"column:productname"`
	ProductNum   int64  `json:"ProductNum"  gorm:"column:productnum" `
	ProductImage string `json:"ProductImage" gorm:"column:productimage" `
	ProductUrl   string `json:"ProductUrl" gorm:"column:producturl"`
}
