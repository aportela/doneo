package rolerepository

type roleDTO struct {
	ID                 string `db:"id"`
	Name               string `db:"name"`
	PermissionsBitmask uint64 `db:"permissions_bitmask"`
}

type searchFilterDTO struct {
	Name *string
}
