package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"myshoppingmall/datamodels"
	"myshoppingmall/encrypt"
	"myshoppingmall/services"
	"net/http"
	"strconv"
)

type UserController struct {
	us *services.UserService
}
func NewUserController(db *gorm.DB)*UserController{
	return &UserController{services.Newuserservice(db)}
}

func (u *UserController) Router(app *gin.Engine) {
	app.GET("/user/register", u.Register)
	app.POST("/user/pregister",u.PRegister)
	app.GET("/user/login",u.Login)
	app.POST("/user/plogin",u.PLogin)
}

func (u *UserController) Register(c *gin.Context)  {
	c.HTML(http.StatusOK,"user/register.html",nil)
}

func (u *UserController) PRegister(c *gin.Context) {
	var (
		nickName = c.PostForm("nickName")
		userName = c.PostForm("userName")
		password = c.PostForm("password")
	)
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}
	_, err := u.us.AddUser(user)
	if err != nil {
		fmt.Printf("err:%v",err)
		return
	}
	c.Redirect(http.StatusMovedPermanently,"login")
}

func (u *UserController) Login(c *gin.Context)  {
	c.HTML(http.StatusOK,"user/login.html",nil)
}

func (u *UserController) PLogin(c *gin.Context) {
	//获取用户提交的表单信息
	var (
		userName = c.PostForm("userName")
		password = c.PostForm("password")
	)
	//验证账号密码正确
	user, isOk := u.us.IsPwdSuccess(userName, password)
	if !isOk {
		//密码错误
		c.HTML(http.StatusOK,"user/login.html",nil)
	}
	//登陆成功
	//写入用户ID到cookie中
	c.SetCookie("uid", strconv.FormatInt(user.ID, 10),
		1000, "/", "localhost", false, true)
	//加密
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//把加密信息写入用户浏览器
	c.SetCookie("sign",uidString,
		3600, "/", "localhost", false, true)
	c.Redirect(http.StatusMovedPermanently,"http://localhost:8080/product/all")
}