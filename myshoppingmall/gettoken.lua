---(r,key1,key2,最大令牌数,放置速率)
redis.replicate_commands()
local nowtimearray=redis.call("time")
--- 获取当前时间，以毫秒为单位
local nowtime=nowtimearray[1]*1000+nowtimearray[2]/1000
redis.call("msetnx",KEYS[1],10,KEYS[2],nowtime)
--- 获取当前桶内令牌数
local currtokens =tonumber(redis.call("get",KEYS[1]))
--- 获取最后一次放入令牌时间
local lastupdatetime =tonumber(redis.call("get",KEYS[2]))
---生成令牌数
local generatecnt=math.floor((nowtime-lastupdatetime)/tonumber(ARGV[2]))
local newtokens=currtokens+generatecnt

if (newtokens>0) then
    if (newtokens>tonumber(ARGV[1]))then
        redis.call("set",KEYS[1],ARGV[1])
    else
        redis.call("set",KEYS[1],newtokens)
    end
    redis.call("DECR", KEYS[1])
    redis.call("set",KEYS[2],nowtime)
    return 1
end
return 0