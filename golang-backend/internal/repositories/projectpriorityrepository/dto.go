package projectpriorityrepository

type ProjectPriorityDTO struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	HexColor string `db:"item_hex_color"`
}

type searchFilterDTO struct {
	Name *string
}
