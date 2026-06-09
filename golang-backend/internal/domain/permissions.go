package domain

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
