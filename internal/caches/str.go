package caches

import "sync"

type StrThreadSafe struct {
	cacheMutex sync.Mutex
	cache      map[string]bool
}

func NewStrThreadSafe() StrThreadSafe {
	return StrThreadSafe{
		cache: make(map[string]bool),
	}
}

func (strCache *StrThreadSafe) Add(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	strCache.cache[str] = true
}

func (strCache *StrThreadSafe) Remove(str string) {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	delete(strCache.cache, str)
}

func (strCache *StrThreadSafe) Has(str string) bool {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	_, hasKey := strCache.cache[str]
	return hasKey
}

func (strCache *StrThreadSafe) Count() int {
	defer strCache.cacheMutex.Unlock()
	strCache.cacheMutex.Lock()

	return len(strCache.cache)
}