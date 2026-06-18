package cache

import (
	"sync"

	"github.com/aportela/doneo/internal/domain"
)

type userPermissions struct {
	GlobalPermissions  domain.Bitmask
	ProjectPermissions map[string]domain.Bitmask
}

type PermissionCache interface {
	GetUser(userID string) (domain.Bitmask, bool)
	SetUser(userID string, permissions domain.Bitmask)

	GetProject(userID, projectID string) (domain.Bitmask, bool)
	SetProject(userID, projectID string, permissions domain.Bitmask)

	DeleteUser(userID string)
	DeleteProject(userID, projectID string)

	Clear()
}

type permissionCache struct {
	mu sync.RWMutex

	permissions map[string]*userPermissions
}

func NewPermissionCache() PermissionCache {
	return &permissionCache{
		permissions: make(map[string]*userPermissions),
	}
}

func (c *permissionCache) GetUser(userID string) (domain.Bitmask, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userPerms, ok := c.permissions[userID]
	if !ok {
		return 0, false
	}

	return userPerms.GlobalPermissions, true
}

func (c *permissionCache) SetUser(userID string, permissions domain.Bitmask) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userPerms, ok := c.permissions[userID]
	if !ok {
		userPerms = &userPermissions{
			ProjectPermissions: make(map[string]domain.Bitmask),
		}
		c.permissions[userID] = userPerms
	}

	userPerms.GlobalPermissions = permissions
}

func (c *permissionCache) GetProject(userID, projectID string) (domain.Bitmask, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userPerms, ok := c.permissions[userID]
	if !ok {
		return 0, false
	}

	perm, ok := userPerms.ProjectPermissions[projectID]
	if !ok {
		return 0, false
	}

	return perm, true
}

func (c *permissionCache) SetProject(userID, projectID string, permissions domain.Bitmask) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userPerms, ok := c.permissions[userID]
	if !ok {
		userPerms = &userPermissions{
			ProjectPermissions: make(map[string]domain.Bitmask),
		}
		c.permissions[userID] = userPerms
	}

	userPerms.ProjectPermissions[projectID] = permissions
}

func (c *permissionCache) DeleteUser(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.permissions, userID)
}

func (c *permissionCache) DeleteProject(userID, projectID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userPerms, ok := c.permissions[userID]
	if !ok {
		return
	}

	delete(userPerms.ProjectPermissions, projectID)
}

func (c *permissionCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.permissions = make(map[string]*userPermissions)
}
