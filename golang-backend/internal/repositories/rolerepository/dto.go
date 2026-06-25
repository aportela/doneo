package rolerepository

type roleBaseDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type roleDTO struct {
	roleBaseDTO
	PermissionsBitmask uint64 `db:"permissions_bitmask"`
}

type searchFilterDTO struct {
	Name *string
}
