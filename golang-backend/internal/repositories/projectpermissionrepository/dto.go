package projectpermissionrepository

type projectPermissionDTO struct {
	ID                     string `db:"id"`
	UserId                 string `db:"user_id"`
	UserName               string `db:"user_name"`
	RoleId                 string `db:"role_id"`
	RoleName               string `db:"role_name"`
	RolePermissionsBitmask uint64 `db:"permissions_bitmask"`
}
