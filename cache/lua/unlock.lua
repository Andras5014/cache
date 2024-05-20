--keys1分布式锁的key
--argc1就是你预期的存在redis的key
if redis.call("GET", KEYS[1]) == ARGV[1] then
    -- 确实是你的锁
    return redis.call("DEL", KEYS[1])
else
    --不是你的锁
    return 0
end