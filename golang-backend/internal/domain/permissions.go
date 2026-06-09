package domain

// user permissions
const (
	UserPermissionAdmin Bitmask = 1 << iota
)

// app (role) permissions
const (
	PermissionUpdateProject Bitmask = 1 << iota
	PermissionDeleteProject
	PermissionViewProject
	PermissionAddTask
	PermissionUpdateTask
	PermissionDeleteTask
	PermissionViewTask
)
