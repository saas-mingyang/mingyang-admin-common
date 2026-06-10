package redis

// ReleaseLockScript 释放锁的 Lua 脚本: 只删除值等于自己 instance 的锁,避免误删别人的
const ReleaseLockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
end
return 0
`
