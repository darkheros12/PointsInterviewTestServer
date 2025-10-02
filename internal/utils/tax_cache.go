package utils

import (
	"PointsInterviewTestServer/internal/models"
	"sync"
)

type item struct {
	brackets []models.TaxBracket
}

type MemoryCache struct {
	// map are not safe for concurrent use (for future scalability)
	mutex sync.RWMutex
	items map[int]item
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[int]item),
	}
}

func (cache *MemoryCache) Get(year int) ([]models.TaxBracket, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	if result, ok := cache.items[year]; ok {
		return result.brackets, true
	}

	return nil, false
}

func (cache *MemoryCache) Set(year int, brackets []models.TaxBracket) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.items[year] = item{
		brackets: brackets,
	}
}
