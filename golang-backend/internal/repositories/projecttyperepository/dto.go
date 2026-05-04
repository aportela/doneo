package projecttyperepository

type projectTypeDTO struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	HexColor string `db:"item_hex_color"`
}

type projectTypeFilterDTO struct {
}
