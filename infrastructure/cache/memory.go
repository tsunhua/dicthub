package cache

import (
	"app/infrastructure/log"
	"fmt"
	"sync"

	"github.com/bluele/gcache"
)

var cache gcache.Cache
var mutex sync.Mutex

func Cache() gcache.Cache {
	return cache
}

func init() {
	mutex.Lock()
	defer mutex.Unlock()

	log.Info("start init gcache")
	cache = gcache.New(1000).LRU().AddedFunc(func(key, value interface{}) {
		log.Debug(fmt.Sprintf("gcache added key:%s", key))
	}).EvictedFunc(func(key, value interface{}) {
		log.Debug(fmt.Sprintf("gcache evicted key:%s", key))
	}).PurgeVisitorFunc(func(key, value interface{}) {
		log.Debug(fmt.Sprintf("gcache purge key:%s", key))
	}).Build()
	log.Info("finish init gcache")
}
