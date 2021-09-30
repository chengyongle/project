package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"myshoppingmall/datamodels"
)

type Userdao struct {
	Db  *gorm.DB
}
func (u *Userdao) SelectByName(userName string) (*datamodels.User, error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	userResult := &datamodels.User{}
	if err:=u.Db.Where("userName=?",userName).First(userResult).Error;
		err==gorm.ErrRecordNotFound{
		return &datamodels.User{},gorm.ErrRecordNotFound
	}
	return  userResult,nil
}

func (u *Userdao) Insert(user *datamodels.User) (int64,  error) {
	err:=u.Db.Create(user).Error
	if err!=nil{
		return  0,err
	}
	var id []int64
	u.Db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return  id[0],err
}

func (u *Userdao) SelectByID(userId int) (*datamodels.User, error) {
	userResult := &datamodels.User{}
	if err:=u.Db.Where("id=?",userId).First(userResult).Error;
		err==gorm.ErrRecordNotFound{
		return &datamodels.User{},gorm.ErrRecordNotFound
	}
	return  userResult,nil
}

