package domain

type Role struct {
	ID                 string
	Name               string
	PermissionsBitmask PermissionsBitmask
}

type SearchRolesFilter struct {
	Name                        *string
	RequiredPermissionsBitmask  *PermissionsBitmask
	ForbiddenPermissionsBitmask *PermissionsBitmask
}
