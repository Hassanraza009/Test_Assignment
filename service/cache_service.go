package service

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type cacheService struct {
	client *cache.Cache
}

type CacheService interface {
	GetValue(key string) (interface{}, bool)
	SetValue(key string, val interface{}, d time.Duration)
}

func NewCacheService() CacheService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &cacheService{
		client: c,
	}
}

func (s *cacheService) GetValue(key string) (interface{}, bool) {
	return s.client.Get(key)
}
func (s *cacheService) SetValue(key string, val interface{}, d time.Duration) {
	s.client.Set(key, val, d)
}
