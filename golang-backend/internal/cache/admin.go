package cache

import "sync"

type AdminCache struct {
	mu    sync.RWMutex
	cache map[string]bool
}

func NewAdminCache() *AdminCache {
	return &AdminCache{
		cache: make(map[string]bool),
	}
}

func (a *AdminCache) SetAdmin(userID string, isAdmin bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[userID] = isAdmin
}

func (a *AdminCache) Delete(userID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.cache, userID)
}

func (a *AdminCache) IsAdmin(userID string, fetchFunc func(string) bool) bool {
	a.mu.RLock()
	val, exists := a.cache[userID]
	a.mu.RUnlock()

	if exists {
		return val
	}

	isAdmin := fetchFunc(userID)
	a.SetAdmin(userID, isAdmin)
	return isAdmin
}
