package redispool

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const(
	//配置连接池参数
	maxidle   = 16
	maxactice = 1000
	idletime  = 300 * time.Second
	redisaddr = "127.0.0.1:6379"
)

func Redispoolinit() *redis.Pool {
	// 建立连接池
	return &redis.Pool{
		MaxIdle:     maxidle,   //最初的连接数量
		MaxActive:   maxactice, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: idletime,  //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", redisaddr)
		},
	}
}
