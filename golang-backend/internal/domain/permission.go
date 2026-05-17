package domain

type PermissionsBitmask uint8

// ProjectView|ProjectUpdate|ProjectDelete|TaskCreate|TaskUpdate|TaskDelete|TaskView
const (
	PermissionCreate PermissionsBitmask = 1 << iota
	PermissionUpdate
	PermissionDelete
	PermissionView
	PermissionList
	PermissionExecute
)

func (p PermissionsBitmask) HasPermission(v PermissionsBitmask) bool {
	return p&v == v
}

func (p *PermissionsBitmask) AddPermission(v PermissionsBitmask) {
	*p |= v
}

func (p *PermissionsBitmask) RemovePermission(v PermissionsBitmask) {
	*p &^= v
}
