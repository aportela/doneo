package domain

type Bitmask uint64

func (p Bitmask) HasFlag(v Bitmask) bool {
	return p&v == v
}

func (p Bitmask) HasAny(v Bitmask) bool {
	return p&v != 0
}

func (p *Bitmask) AddFlag(v Bitmask) {
	*p |= v
}

func (p *Bitmask) RemoveFlag(v Bitmask) {
	*p &^= v
}

func (p *Bitmask) ToggleFlag(v Bitmask) {
	*p ^= v
}

func (p *Bitmask) Clear() {
	*p = 0
}

// user permissions
const (
	UserPermissionAdmin Bitmask = 1 << iota
)

// app permissions
const (
	PermissionUpdateProject Bitmask = 1 << iota
	PermissionDeleteProject
	PermissionViewProject
	PermissionAddTask
	PermissionUpdateTask
	PermissionDeleteTask
	PermissionViewTask
)

// app permissions
const (
	PermissionCreate Bitmask = 1 << iota
	PermissionUpdate
	PermissionDelete
	PermissionView
	PermissionList
	PermissionExecute
)
