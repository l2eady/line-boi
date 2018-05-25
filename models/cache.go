package models

import "github.com/bluele/gcache"

type CacheService struct {
	Gcache gcache.Cache
}
