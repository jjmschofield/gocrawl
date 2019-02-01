package caches

import "sync"

type StrLocking struct {
	cacheMutex sync.Mutex
	cache      map[string]interface{}
}

func NewStrBlocking() StrLocking {
	return StrLocking{
		cache: make(map[string]interface{}),
	}
}

func (strCache *StrLocking) Has(str string) bool {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	_, hasKey := strCache.cache[str]
	return hasKey
}

func (strCache *StrLocking) Add(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	strCache.cache[str] = nil
}

func (strCache *StrLocking) Remove(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	delete(strCache.cache, str)
}

func (strCache *StrLocking) Count() int {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	return len(strCache.cache)
}