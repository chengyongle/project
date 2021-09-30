package datamodels

type Message struct {
	ProductID int64
	UserID    int64
}
//初始化消息结构体
func NewMessage(userId int64,productId int64) *Message  {
	return &Message{UserID:userId,ProductID:productId}
}