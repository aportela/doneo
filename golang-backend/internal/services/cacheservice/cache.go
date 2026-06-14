package cacheservice

import (
	"sync"

	"github.com/aportela/doneo/internal/domain"
)

type UserPermissionCacheKey struct {
	UserID string
}

type UserPermissionCache interface {
	Get(userID string) (domain.Bitmask, bool)
	Set(userID string, permissionBitmask domain.Bitmask)
	Delete(userID string)
	Clear()
}

type userPermissionCache struct {
	mu          sync.RWMutex
	permissions map[UserPermissionCacheKey]domain.Bitmask
}

func NewUserPermissionCache() UserPermissionCache {
	return &userPermissionCache{
		permissions: make(map[UserPermissionCacheKey]domain.Bitmask),
	}
}

func (service *userPermissionCache) Get(userID string) (domain.Bitmask, bool) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	permissionBitmask, ok := service.permissions[UserPermissionCacheKey{
		UserID: userID,
	}]

	return permissionBitmask, ok
}

func (service *userPermissionCache) Set(userID string, permissionBitmask domain.Bitmask) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions[UserPermissionCacheKey{
		UserID: userID,
	}] = permissionBitmask
}

func (service *userPermissionCache) Delete(userID string) {
	service.mu.Lock()
	defer service.mu.Unlock()

	delete(service.permissions, UserPermissionCacheKey{
		UserID: userID,
	})
}

func (service *userPermissionCache) Clear() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions = make(map[UserPermissionCacheKey]domain.Bitmask)
}

type ProjectPermissionCacheKey struct {
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
	permissions map[ProjectPermissionCacheKey]domain.Bitmask
}

func NewProjectPermissionCache() ProjectPermissionCache {
	return &projectPermissionCache{
		permissions: make(map[ProjectPermissionCacheKey]domain.Bitmask),
	}
}

func (service *projectPermissionCache) Get(userID string, projectID string) (domain.Bitmask, bool) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	permissionBitmask, ok := service.permissions[ProjectPermissionCacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}]

	return permissionBitmask, ok
}

func (service *projectPermissionCache) Set(userID string, projectID string, permissionBitmask domain.Bitmask) {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions[ProjectPermissionCacheKey{
		UserID:    userID,
		ProjectID: projectID,
	}] = permissionBitmask
}

func (service *projectPermissionCache) Delete(userID string, projectID string) {
	service.mu.Lock()
	defer service.mu.Unlock()

	delete(service.permissions, ProjectPermissionCacheKey{
		UserID:    userID,
		ProjectID: projectID,
	})
}

func (service *projectPermissionCache) Clear() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.permissions = make(map[ProjectPermissionCacheKey]domain.Bitmask)
}
