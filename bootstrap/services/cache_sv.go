package services

import (
	"github.com/ggoop/mdf/gmap"
	"strings"
)

type ICacheSv interface {
	Push(key string, value interface{})
	Has(key string) bool
	Remove(key string)
}
type cacheSvImpl struct {
	cache gmap.ConcurrentMap
}

var cacheSvInstance ICacheSv = newCacheSv()

func CacheSv() ICacheSv {
	return cacheSvInstance
}
func newCacheSv() *cacheSvImpl {
	return &cacheSvImpl{cache: gmap.New()}
}

func (s *cacheSvImpl) Push(key string, value interface{}) {
	s.cache.Set(strings.ToLower(key), value)
}
func (s *cacheSvImpl) Get(key string) interface{} {
	if val, ok := s.cache.Get(strings.ToLower(key)); ok {
		return val
	}
	return nil
}
func (s *cacheSvImpl) Has(key string) bool {
	return s.cache.Has(strings.ToLower(key))
}
func (s *cacheSvImpl) Remove(key string) {
	s.cache.Remove(strings.ToLower(key))
}
