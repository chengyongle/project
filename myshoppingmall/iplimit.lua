local notexists = redis.call("set", KEYS[1], 1, "NX", "EX", tonumber(ARGV[2]))
if (notexists) then
    return 1
end
local current = tonumber(redis.call("get", KEYS[1]))
if (current == nil) then
    redis.call("incr", KEYS[1])
    redis.call("expire", KEYS[1], tonumber(ARGV[2]))
    return 1
end
if (current >= tonumber(ARGV[1])) then
    return 2
end
redis.call("incr", KEYS[1])
return 1