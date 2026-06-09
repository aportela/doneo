package domain

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
