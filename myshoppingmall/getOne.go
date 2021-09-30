package main

import (
	"fmt"
	"myshoppingmall/tokenbucket"
	"net/http"
	"sync"
)

//配置参数
const (
	maxtokens = 10 //最大令牌数
	rate      = 20 //令牌生成速率：几毫秒一个
)

//创建全局变量
var (
	tb *tokenbucket.TokenBucket//令牌桶
	sum        int64 = 0 //已抢数量
	productNum int64 = 100//预存商品数量
	mutex sync.Mutex //互斥锁
)

//获取秒杀商品
func GetOneProduct() bool {
	//加锁
	mutex.Lock()
	defer mutex.Unlock()
	if tb.GetBucket() {
		fmt.Println("获取令牌成功")
		if sum < productNum {
			sum += 1
			fmt.Println(sum)
			return true
		}
	}
	fmt.Println("获取令牌失败")
	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		//抢到
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	tb = tokenbucket.NewTokenBucket("nowtokens", "lastupdatetime", maxtokens, rate)
	http.HandleFunc("/getOne", GetProduct)
	http.ListenAndServe(":8083", nil)
}
