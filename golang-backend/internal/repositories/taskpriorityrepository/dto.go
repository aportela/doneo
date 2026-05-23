package taskpriorityrepository

type taskPriorityDTO struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	HexColor string `db:"item_hex_color"`
}

type searchFilterDTO struct {
	Name *string
}
