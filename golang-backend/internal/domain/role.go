package domain

const (
	PermissionUpdateProject Bitmask = 1 << iota
	PermissionDeleteProject
	PermissionViewProject
	PermissionAddTask
	PermissionUpdateTask
	PermissionDeleteTask
	PermissionViewTask
)

type RoleBase struct {
	ID   string
	Name string
}

type Role struct {
	RoleBase
	PermissionsBitmask Bitmask
}

type SearchRolesFilter struct {
	Name                        *string
	RequiredPermissionsBitmask  *Bitmask
	ForbiddenPermissionsBitmask *Bitmask
}
