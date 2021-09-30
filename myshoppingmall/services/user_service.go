package services

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"myshoppingmall/dao"
	"myshoppingmall/datamodels"
)

type UserService struct {
	userdao dao.Userdao
}
func Newuserservice(db *gorm.DB) *UserService{
	ud:=dao.Userdao{db}
	return  &UserService{ud}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (*datamodels.User,bool) {
	user, err := u.userdao.SelectByName(userName)
	if err != nil {
		//用户不存在
		fmt.Printf("err:%v",err)
		return &datamodels.User{},false
	}
	isOk,_:= ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return user,isOk
}

func (u *UserService) AddUser(user *datamodels.User) (int64, error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return 0, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.userdao.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (bool, error) {
	if err:= bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}
