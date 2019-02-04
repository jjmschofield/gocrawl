package caches

import "sync"

type LockingStr struct {
	cacheMutex sync.Mutex
	cache      map[string]interface{}
}

func NewLockingStr() LockingStr {
	return LockingStr{
		cache: make(map[string]interface{}),
	}
}

func (strCache *LockingStr) Add(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	strCache.cache[str] = nil
}

func (strCache *LockingStr) Remove(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	delete(strCache.cache, str)
}

func (strCache *LockingStr) Has(str string) bool {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	_, hasKey := strCache.cache[str]
	return hasKey
}

func (strCache *LockingStr) Count() int {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	return len(strCache.cache)
}