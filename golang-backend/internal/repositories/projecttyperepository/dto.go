package projecttyperepository

type projectTypeDTO struct {
	ID                string `db:"id"`
	Name              string `db:"name"`
	HexColor          string `db:"item_hex_color"`
	WorkspaceId       string `db:"workspace_id"`
	WorkspaceName     string `db:"workspace_name"`
	WorkspaceHexColor string `db:"workspace_hex_color"`
}

type projectTypeFilterDTO struct {
	WorkspaceId string
}
