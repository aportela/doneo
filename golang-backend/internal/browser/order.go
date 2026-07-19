package browser

type Order struct {
	Field     string
	Direction string
}

func (p Order) hasAscendingSort() bool {
	return p.Direction == "ASC"
}

func (p Order) hasDescendingSort() bool {
	return p.Direction == "DESC"
}
