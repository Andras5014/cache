val = redis.call('get',keys[1])
if val ==false then
    -- key not exist
    return redis.call('set',keys[1],argc[1],'ex',argv[2])
elseif val == argv[1] then
    -- last lock success

    redis .call('expire',keys[1],argv[2])
    return 'OK'
else
    --lock exist
    return ''
end