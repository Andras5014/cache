if redis.call("GET", KEYS[1]) == ARGV[1] then
    -- 确实是你的锁
    return redis.call("expire", KEYS[1],argv[2])
else
    --不是你的锁
    return 0
end