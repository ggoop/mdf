package service

import (
	"github.com/ggoop/mdf/gmap"
	"strings"
	"sync"
)

type CacheSv struct {
	cache gmap.ConcurrentMap
}

var _cacheInstance *CacheSv
var _cacheDefaultMu sync.Once

func NewCacheSv() *CacheSv {
	_cacheDefaultMu.Do(func() {
		_cacheInstance = &CacheSv{cache: gmap.New()}
	})
	return _cacheInstance
}
func (s *CacheSv) Push(key string, value interface{}) {
	s.cache.Set(strings.ToLower(key), value)
}
func (s *CacheSv) Get(key string) interface{} {
	if val, ok := s.cache.Get(strings.ToLower(key)); ok {
		return val
	}
	return nil
}
func (s *CacheSv) Has(key string) bool {
	return s.cache.Has(strings.ToLower(key))
}
func (s *CacheSv) Remove(key string) {
	s.cache.Remove(strings.ToLower(key))
}
