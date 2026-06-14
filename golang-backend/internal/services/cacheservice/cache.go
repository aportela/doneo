package cacheservice

import (
	"sync"

	"github.com/aportela/doneo/internal/domain"
)

type CacheKey struct {
	UserID    string
	ProjectID string
}

type ProjectPermissionCache interface {
	Get(userID string, projectID string) (domain.Bitmask, bool)
	Set(userID string, projectID string, permissionBitmask domain.Bitmask)
	Delete(userID string, projectID string)
	Clear()
}

type projectPermissionCache struct {
	mu          sync.RWMutex
	permissions map[CacheKey]domain.Bitmask
}

func NewProjectPermissionCache() ProjectPermissionCache {
	return &projectPermissionCache{
		permissions: make(map[CacheKey]domain.Bitmask),
	}
}

func (service *projectPermissionCache) Get(userID string, projectID string) (domain.Bitmask, bool) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	permissionBitmask, ok := service.permissions[CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}]

	return permissionBitmask, ok
}

func (service *projectPermissionCache) Set(userID string, projectID string, permissionBitmask domain.Bitmask) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions[CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}] = permissionBitmask
}

func (service *projectPermissionCache) Delete(userID string, projectID string) {
	service.mu.Lock()
	defer service.mu.Unlock()

	delete(service.permissions, CacheKey{
		UserID:    userID,
		ProjectID: projectID,
	})
}

func (service *projectPermissionCache) Clear() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions = make(map[CacheKey]domain.Bitmask)
}
