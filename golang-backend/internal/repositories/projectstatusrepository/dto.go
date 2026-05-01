package projectstatusrepository

type projectStatusDTO struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Index    int    `db:"item_index"`
	HexColor string `db:"item_hex_color"`
}
