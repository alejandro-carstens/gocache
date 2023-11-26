package gocache

const (
	redisLuaReleaseLockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`
	redisLuaAddScript = `
return redis.call('exists',KEYS[1])<1 and redis.call('setex',KEYS[1],ARGV[2],ARGV[1])
`
	redisLuaExpireLockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
	if redis.call("del",KEYS[1]) == 0 then
		return 0;
	end
    if ARGV[2] == "0" or redis.call('setex',KEYS[1],ARGV[2],ARGV[1]) then
		return 1;
	end
end

return 0;
`
)
