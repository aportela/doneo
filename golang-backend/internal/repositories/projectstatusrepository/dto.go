package projectstatusrepository

type projectStatusDTO struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	HexColor string `db:"item_hex_color"`
	Index    int    `db:"item_index"`
}

type searchFilterDTO struct {
	Name *string
}
