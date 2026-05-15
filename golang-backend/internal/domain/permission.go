package domain

type PermissionBitmask uint8

const (
	PermissionCreate PermissionBitmask = 1 << iota
	PermissionUpdate
	PermissionDelete
	PermissionView
	PermissionList
	PermissionExecute
)

func (p PermissionBitmask) HasPermission(v PermissionBitmask) bool {
	return p&v == v
}

func (p *PermissionBitmask) AddPermission(v PermissionBitmask) {
	*p |= v
}

func (p *PermissionBitmask) RemovePermission(v PermissionBitmask) {
	*p &^= v
}
