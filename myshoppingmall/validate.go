package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"myshoppingmall/common"
	"myshoppingmall/datamodels"
	"myshoppingmall/encrypt"
	"myshoppingmall/rabbitmq"
	"myshoppingmall/redispool"
	"net/http"
	"net/url"
	"strconv"
)
//配置参数
const (
	GetOneIp = "172.19.0.7"
	GetOnePort = "8083"
	port = "8080"
	maxtime=60 //过期时间 单位秒
	maxcnt=100 //最大访问数
)
//创建全局变量
var (
	localHost = ""
	hostArray= []string{"172.19.0.6","172.19.0.8"} //设置集群地址
	hashConsistent *common.Consistent
	rabbitMqValidate *rabbitmq.RabbitMQ
	accessControl *AccessControl
)
//访问控制
type AccessControl struct {
	lua *redis.Script
	redisClient *redis.Pool //连接池
}
func NewaccessControl(rs *redis.Script)*AccessControl{
	return &AccessControl{rs,redispool.Redispoolinit()}
}

//分布式验证
func (m *AccessControl) GetDistributedRight(req *http.Request) bool{
	//获取用户UID
	uid,err:=req.Cookie("uid")
	if err !=nil {
		fmt.Println("未登录")
		return false
	}
	//采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest,err:=hashConsistent.Get(uid.Value)
	if err !=nil {
		fmt.Println(err)
		return false
	}
	//获取用户IP
	uip:=req.RemoteAddr
	fmt.Println("RemoteAddr:", req.RemoteAddr)
	uinfo:=uid.Value+":"+uip
	//判断是否为本机
	if hostRequest==localHost{
		//执行本机数据读取和校验
		return m.Localvalidate(uinfo)
	}else {
		//不是本机充当代理访问数据返回结果
		return Proxyvalidate(hostRequest,req)
	}

}
//本地验证
func (m *AccessControl) Localvalidate(uinfo string) bool {
	//从连接池中取一个连接
	rdb:=m.redisClient.Get()
	defer rdb.Close()
	res,err:=m.lua.Do(rdb,uinfo,maxcnt,maxtime)
	if err!=nil{
		fmt.Println("false10")
		fmt.Println(err)
		return  false
	}
	//验证成功
	if res==int64(1){
		return true
	}
	fmt.Println("请求过多")
	return  false
}
//发送代理请求
func Sendreq(hostUrl string,request *http.Request)(response *http.Response,body []byte,err error){
	//获取Uid和sign
	uidPre,err := request.Cookie("uid")
	if err !=nil {
		return
	}
	uidSign,err:=request.Cookie("sign")
	if err !=nil {
		return
	}
	//模拟接口访问
	client :=&http.Client{}
	req,err:=http.NewRequest("GET",hostUrl,nil)
	if err !=nil {
		return
	}
	//手动指定
	cookieUid :=&http.Cookie{Name:"uid",Value:uidPre.Value,Path:"/"}
	cookieSign :=&http.Cookie{Name:"sign",Value:uidSign.Value,Path:"/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	//获取返回结果
	response,err =client.Do(req)
	if err !=nil {
		return
	}
	body,err=ioutil.ReadAll(response.Body)
	return
}
//获取其它节点处理结果
func Proxyvalidate(host string,request *http.Request) bool  {
	hostUrl:="http://"+host+":"+port+"/checkRight"
	response,body,err:=Sendreq(hostUrl,request)
	if err !=nil {
		return false
	}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}
//验证
func CheckRight(w http.ResponseWriter,r *http.Request)  {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}
//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	queryForm,err:=url.ParseQuery(r.URL.RawQuery)
	if err !=nil || len(queryForm["productID"])<=0 {
		w.Write([]byte("false9"))
		return
	}
	productString :=queryForm["productID"][0]
	fmt.Printf("productid:%v\n",productString)
	//获取用户cookie
	userCookie,err:=r.Cookie("uid")
	if err !=nil {
		w.Write([]byte("false8"))
		return
	}
	//分布式权限验证
	right:=accessControl.GetDistributedRight(r)
	if right == false{
		w.Write([]byte("false7"))
		return
	}
	//获取数量控制权限，防止超卖
	hostUrl :="http://"+GetOneIp+":"+GetOnePort+"/getOne"
	responseValidate,validateBody,err:=Sendreq(hostUrl,r)
	if err !=nil {
		w.Write([]byte("false0"))
		return
	}
	fmt.Println(responseValidate.StatusCode)
	if responseValidate.StatusCode==200{
		fmt.Println(string(validateBody))
		if string(validateBody)=="true"{
			//下单
			//获取商品ID
			productID,err:=strconv.ParseInt(productString,10,64)
			if err!=nil{
				w.Write([]byte("false1"))
				return
			}
			//获取用户ID
			userID,err := strconv.ParseInt(userCookie.Value,10,64)
			if err !=nil {
				w.Write([]byte("false2"))
				return
			}
			//创建消息体
			message:=datamodels.NewMessage(userID,productID)
			byteMessage,err :=json.Marshal(&message)
			//生产消息
			err=rabbitMqValidate.PublishSimple(string(byteMessage))
			if err!=nil{
				fmt.Println(err)
				w.Write([]byte("false3"))
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false4"))
	return
}

//统一验证拦截器
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

//身份校验函数
func CheckUserInfo(r *http.Request) error {
	//获取Uid，cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		fmt.Println("用户UID Cookie 获取失败！")
		fmt.Printf("err:%v",err)
		return err
	}
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		fmt.Println("用户加密串 Cookie 获取失败！")
		return err
	}
	//对信息进行解密
	signc,_:=url.PathUnescape(signCookie.Value)
	signByte, err := encrypt.DePwdCode(signc)
	if err != nil {
		fmt.Println("加密串已被篡改！")
		return err
	}
	if checkInfo(uidCookie.Value, string(signByte)) {
		fmt.Println("身份校验成功！")
		return nil
	}
	return errors.New("身份校验失败！")
}

//判断相等方法
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}
//初始化一些配置
func initvd(){
	//初始化lua脚本
	script, err := ioutil.ReadFile("./iplimit.lua")
	if err != nil {
		fmt.Println("Script read error", err)
		return
	}
	l:=redis.NewScript(1,string(script))
	//初始化访问控制器
	accessControl=NewaccessControl(l)
	//负载均衡器设置
	hashConsistent =common.NewConsistent()
	//一致性hash算法，添加节点
	for _,v :=range hostArray {
		hashConsistent.Add(v)
	}
	//获取本机IP
	localIp,err:=common.GetIntranceIp()
	if err!=nil {
		fmt.Println(err)
	}
	localHost=localIp
	fmt.Println(localHost)
}
func main() {
	initvd()
	//初始化rmq
	rabbitMqValidate=rabbitmq.NewRabbitMQSimple("myproduct")
	defer rabbitMqValidate.Destory()
	//过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight",Auth)
	//启动拦截器
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight",filter.Handle(CheckRight))
	//启动服务
	http.ListenAndServe(":8080", nil)
}

