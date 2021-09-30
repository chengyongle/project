package tokenbucket

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"myshoppingmall/redispool"
)


//定义令牌桶结构体
type TokenBucket struct {
	key1      string //redis中当前桶内令牌数key
	key2      string	//redis中当上次更新时间key
	maxtokens int //最大令牌数
	rate  int //令牌生成速率：几毫秒一个
	lua *redis.Script //lua脚本
	redisClient *redis.Pool //连接池
}

//初始化令牌桶
func NewTokenBucket(k1,k2 string,mt,rt int) *TokenBucket{
	//初始化lua脚本
	script, err := ioutil.ReadFile("./gettoken.lua")
	if err != nil {
		fmt.Println("Script read error", err)
		return nil
	}
	l:=redis.NewScript(2,string(script))
	return &TokenBucket{
		key1:      k1,
		key2:      k2,
		maxtokens: mt,
		rate:      rt,
		lua:       l,
		redisClient: redispool.Redispoolinit(),
	}
}


//获取令牌
func (t *TokenBucket) GetBucket() bool{
	rconn:=t.redisClient.Get()
	res,err:=t.lua.Do(rconn,t.key1,t.key2,t.maxtokens,t.rate)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	if res!=int64(1){
		return false
	}
	return true
}
