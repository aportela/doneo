package rolerepository

type RoleDTO struct {
	ID                 string `db:"id"`
	Name               string `db:"name"`
	PermissionsBitmask uint8  `db:"permission_bitmask"`
}

type SearchRolesFilterDTO struct {
	Name                        *string
	RequiredPermissionsBitmask  *uint8
	ForbiddenPermissionsBitmask *uint8
}
