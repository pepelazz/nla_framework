package cacheUtil

import (
	"time"
	"github.com/bluele/gcache"
	"fmt"
)

type cacheType struct {
	Data        interface{}
	ExpiredTime time.Time
}

var (
	memCacheMap = map[string]cacheType{}
	gc          = gcache.New(40).
		LRU().
		Build()
)

func MemCacheGet(key string) interface{} {
	if res, ok := memCacheMap[key]; ok {
		// проверяем не истекло ли время кэша
		if res.ExpiredTime.After(time.Now()) {
			return res.Data
		} else {
			delete(memCacheMap, key)
		}
	}
	return nil
}

func MemCachePut(key string, duration int, data interface{}) {
	memCacheMap[key] = cacheType{data, time.Now().Add(time.Duration(duration) * time.Second)}
}

func MemCacheClear(key string) {
	delete(memCacheMap, key)
}

func GoCacheSet(key, value interface{}, t time.Duration) {
	gc.SetWithExpire(key, value, t)
}

func GoCacheGet(key interface{}) (interface{}, error) {
	return gc.Get(key)
}

func GoCacheRemove(key interface{}) bool {
	return gc.Remove(key)
}

// ключ для хранении в кэше данных о польззователе
func GetCacheKeyUser(userId int64) string {
	return fmt.Sprintf("user_id_%v", userId)
}

// ключ для хранении в кэше данных о пользователе по токену
func GetCacheKeyUserToken(token string) string {
	return fmt.Sprintf("token_%v", token)
}

func UserRemoveByToken(token string) bool {
	return GoCacheRemove(GetCacheKeyUserToken(token))
}
