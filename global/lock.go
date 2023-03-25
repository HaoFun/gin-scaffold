package global

import (
	"context"
	"gin-scaffold/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

type Interface interface {
	Get() bool
	Block(seconds int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string
	owner   string
	seconds int64
}

const releaseLockLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

func Lock(name string, seconds int64) Interface {
	return &lock{
		context: context.Background(),
		name:    name,
		owner:   utils.RandString(16),
		seconds: seconds,
	}
}

func (l *lock) Get() bool {
	return App.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// Block for a period of time and attempt to acquire a lock.
func (l *lock) Block(seconds int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-seconds >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// Unlock
func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaScript)
	result := luaScript.Run(l.context, App.Redis, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

// Force unlock
func (l *lock) ForceRelease() {
	App.Redis.Del(l.context, l.name).Val()
}
