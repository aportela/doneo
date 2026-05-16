package rolerepository

type RoleDTO struct {
	ID                string `db:"id"`
	Name              string `db:"name"`
	PermissionBitmask uint8  `db:"permission_bitmask"`
}

type SearchRolesFilterDTO struct {
	Name              *string
	PermissionBitmask *uint8
}
