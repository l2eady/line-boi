package cache

import (
	"fmt"
	"line-boi/models"

	"github.com/bluele/gcache"
)

func NewCache(size int) *models.CacheService {
	return &models.CacheService{
		Gcache: gcache.New(size).
			AddedFunc(func(key, value interface{}) {
				fmt.Println("added key:", key)
			}).Build(),
	}

}
