package rolerepository

type RoleDTO struct {
	ID                 string `db:"id"`
	Name               string `db:"name"`
	PermissionsBitmask uint64 `db:"permissions_bitmask"`
}

type SearchRolesFilterDTO struct {
	Name *string
}
