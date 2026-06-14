package cacheservice

import (
	"sync"

	"github.com/aportela/doneo/internal/domain"
)

type CacheKey struct {
	UserID    string
	ProjectID string
}

type PermissionCache interface {
	Get(userID string, projectID string) (domain.Bitmask, bool)
	Set(userID string, projectID string, perms domain.Bitmask)
	Delete(userID string, projectID string)
	Clear()
}

type permissionCache struct {
	mu    sync.RWMutex
	perms map[CacheKey]domain.Bitmask
}

func NewPermissionCache() PermissionCache {
	return &permissionCache{
		perms: make(map[CacheKey]domain.Bitmask),
	}
}

func (service *permissionCache) Get(userID string, projectID string) (domain.Bitmask, bool) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	perms, ok := service.perms[CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}]

	return perms, ok
}

func (service *permissionCache) Set(userID string, projectID string, perms domain.Bitmask) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.perms[CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}] = perms
}

func (service *permissionCache) Delete(userID string, projectID string) {
	service.mu.Lock()
	defer service.mu.Unlock()

	delete(service.perms, CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	})
}

func (service *permissionCache) Clear() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.perms = make(map[CacheKey]domain.Bitmask)
}
